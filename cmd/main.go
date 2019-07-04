/*******************************************************************************
 * Copyright (c) 2019 Red Hat Inc
 *
 * See the NOTICE file(s) distributed with this work for additional
 * information regarding copyright ownership.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Eclipse Public License 2.0 which is available at
 * http://www.eclipse.org/legal/epl-2.0
 *
 * SPDX-License-Identifier: EPL-2.0
 *******************************************************************************/

package main

import (
	"crypto/tls"
	"log"

	"github.com/spf13/cobra"
)

var insecure bool
var contentType string = "text/plain"
var processCommands bool = false
var ttd uint32 = 0
var qos uint8 = 0

func createTlsConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: insecure,
	}
}

func main() {

	cmdConsume := &cobra.Command{
		Use:   "consume [telemetry|event] [message endpoint uri] [tenant]",
		Short: "Consume and print messages",
		Long:  `Consume messages from the endpoint and print it on the console.`,
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if err := consume(args[0], args[1], args[2]); err != nil {
				log.Fatal("Failed to consume messages:", err)
			}
		},
	}

	cmdPublish := &cobra.Command{
		Use:   "publish",
		Short: "Publish messages",
	}

	cmdPublishHttp := &cobra.Command{
		Use:   "http [telemetry|event] [http endpoint uri] [tenant] [deviceId] [authId] [password] [payload]",
		Short: "Publish via HTTP",
		Args:  cobra.ExactArgs(7),
		Run: func(cmd *cobra.Command, args []string) {
			if err := publishHttp(args[0], args[1], args[2], args[3], args[4], args[5], contentType, args[6]); err != nil {
				log.Fatal("Failed to publish via HTTP:", err)
			}
		},
	}

	cmdPublish.AddCommand(cmdPublishHttp)

	// publish flags

	cmdPublish.PersistentFlags().StringVar(&contentType, "content-type", "text/plain", "content type")

	// publish http flags

	cmdPublishHttp.Flags().Uint32VarP(&ttd, "ttd", "t", 0, "Wait for command")
	cmdPublishHttp.Flags().Uint8VarP(&qos, "qos", "q", 0, "Quality of service")

	// consume flags

	cmdConsume.Flags().BoolVarP(&processCommands, "command", "c", false, "Enable commands")

	// root command

	var rootCmd = &cobra.Command{Use: "hot"}
	rootCmd.AddCommand(cmdConsume, cmdPublish)

	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "Skip TLS validation")

	if err := rootCmd.Execute(); err != nil {
		println(err.Error())
	}
}
