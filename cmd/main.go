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

 //./hot consume telemetry amqps://messaging-8lxzny44dx-enmasse-infra.apps.astoycos-ocp.shiftstack.com:443 myapp.iot --tlsConfig=2 --tlsPath=tls.crt --clientUsername=consumer --clientPassword=foobar

package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"errors"
	"log"
	"strconv"

	"github.com/ctron/hot/pkg/encoding"

	"github.com/spf13/cobra"
)
var tlsConn bool
var insecure bool
var tlsPath string = ""
var clientUsername string = ""
var clientPassword string = ""
var contentTypeFlag string = "text/plain"
var commandReader string = ""
var processCommands bool = false
var ttd uint32 = 0
var qos uint8 = 0

func createTlsConfig() *tls.Config {
	//Insecure TLS 
	if (tlsConn && insecure){
		return &tls.Config{
			InsecureSkipVerify:true,
		}
	//Secure TLS
	} else{
		caCert, err := ioutil.ReadFile(tlsPath)   	
			if err != nil {
				log.Fatal(err)
			}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	
		return &tls.Config{
			RootCAs: caCertPool,
		}
	}
}

func getEncoder() encoding.PayloadEncoder {
	return encoding.CreateEncoder(contentTypeFlag)
}

func main() {

	cmdConsume := &cobra.Command{
		Use:   "consume [telemetry|event] [message endpoint uri] [tenant]",
		Short: "Consume and print messages",
		Long:  `Consume messages from the endpoint and print it on the console.`,
		Args: func(cmd *cobra.Command, args []string) error { 
			cobra.ExactArgs(3)
			if len(args) != 3{
				return errors.New("Wrong number of Input arguments expected 3 got "+strconv.Itoa(len(args)))
			}
			if (tlsConn == false && insecure == true) {
				return errors.New("Cannot have an insecure TLS connection without TLS enabled")
			}
			return nil; 
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := consume(args[0], args[1], args[2]); err != nil {
				log.Fatal("Failed to consume messages: ", err)
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
			if err := publishHttp(HttpPublishInformation{
				CommonPublishInformation: CommonPublishInformation{
					MessageType:      args[0],
					URI:              args[1],
					Tenant:           args[2],
					DeviceId:         args[3],
					AuthenticationId: args[4],
					Password:         args[5],
				},
				QoS: qos,
			}, getEncoder(), args[6]); err != nil {
				log.Fatal("Failed to publish via HTTP:", err)
			}
		},
	}

	cmdPublishMqtt := &cobra.Command{
		Use:   "mqtt [telemetry|event] [mqtt endpoint uri] [tenant] [deviceId] [authId] [password] [payload]",
		Short: "Publish via MQTT",
		Args:  cobra.ExactArgs(7),
		Run: func(cmd *cobra.Command, args []string) {
			if err := publishMqtt(MqttPublishInformation{
				CommonPublishInformation: CommonPublishInformation{
					MessageType:      args[0],
					URI:              args[1],
					Tenant:           args[2],
					DeviceId:         args[3],
					AuthenticationId: args[4],
					Password:         args[5],
				},
				QoS: qos,
			}, getEncoder(), args[6]); err != nil {
				log.Fatal("Failed to publish via MQTT:", err)
			}
		},
	}

	cmdPublish.AddCommand(cmdPublishHttp)
	cmdPublish.AddCommand(cmdPublishMqtt)

	// publish flags

	// publish http flags

	cmdPublishHttp.Flags().Uint32VarP(&ttd, "ttd", "c", 0, "Wait for command")
	cmdPublishHttp.Flags().Uint8VarP(&qos, "qos", "q", 0, "Quality of service")

	// publish mqtt flags

	cmdPublishMqtt.Flags().Uint8VarP(&qos, "qos", "q", 0, "Quality of service")

	// consume flags

	cmdConsume.Flags().BoolVarP(&processCommands, "command", "c", false, "Enable commands")
	cmdConsume.Flags().StringVarP(&commandReader, "reader", "r", "prefill", "Command reader type (possible values: [ondemand, prefill]")
	cmdConsume.Flags().StringVarP(&clientUsername,"clientUsername","u","","Tenant username")
	cmdConsume.Flags().StringVarP(&clientPassword,"clientPassword","p","","Tenant password")
	// root command

	var cmdRoot = &cobra.Command{Use: "hot"}
	cmdRoot.AddCommand(cmdConsume, cmdPublish)

	cmdRoot.PersistentFlags().StringVarP(&contentTypeFlag, "content-type", "t", "text/plain", "Content type of the payload, may be a MIME type or 'hex'")
	cmdRoot.PersistentFlags().BoolVarP(&tlsConn,"tlsConn","T",false,"Set to true to enable a secure TLS connection")
	cmdRoot.PersistentFlags().BoolVarP(&insecure,"insecure","I",false,"Set to true to enable insecure TLS connection")
	cmdRoot.PersistentFlags().StringVarP(&tlsPath,"tlsPath","P","","Absolute path to cert file")

	if err := cmdRoot.Execute(); err != nil {
		println(err.Error())
	}
}
