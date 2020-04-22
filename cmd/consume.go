/*******************************************************************************
 * Copyright (c) 2019, 2020 Red Hat Inc
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
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ctron/hot/pkg/command"

	"github.com/google/uuid"

	"github.com/ctron/hot/pkg/utils"
	"pack.ag/amqp"
)

func createCommandReader(commandName string) (command.Reader, error) {
	lower := strings.ToLower(commandReader)
	if lower == "prefill" {
		return &command.PreFillReader{
			Prompt:  os.Stdout,
			Stream:  os.Stdin,
			Encoder: getEncoder(),
		}, nil
	} else if lower == "ondemand" {
		return &command.OnDemandReader{
			Prompt:  os.Stdout,
			Stream:  os.Stdin,
			Encoder: getEncoder(),
		}, nil
	} else if strings.HasPrefix(lower, "static:") {
		return command.NewStatic(commandName, getEncoder(), commandReader[len("static:"):])
	} else {
		return nil, fmt.Errorf("unknown command reader: %s", commandReader)
	}
}

type cleanup func()

func createReceiver(session *amqp.Session, tenant string, messageType string) (*amqp.Receiver, cleanup, error) {
	receiver, err := session.NewReceiver(
		amqp.LinkSourceAddress(fmt.Sprintf("%s/%s", messageType, tenant)),
		amqp.LinkCredit(10),
	)
	if err != nil {
		return nil, nil, err
	}
	return receiver, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		if err := receiver.Close(ctx); err != nil {
			log.Fatal("Failed to close receiver: ", err)
		}
		cancel()
	}, nil
}

func consume(uri string, tenant string) error {

	if strings.HasPrefix(uri, "amqps:") && disableTlsNegotiation {
		return fmt.Errorf("TLS negotiation is explicitly disabled, but URI indicates TLS: %s", uri)
	}

	fmt.Printf("Consuming from %s for tenant %s ...", uri, tenant)
	fmt.Println()

	opts := make([]amqp.ConnOption, 0)

	opts = append(opts, amqp.ConnTLSConfig(createTlsConfig()))

	// Enable Client credentials if available
	if clientUsername != "" && clientPassword != "" {
		opts = append(opts, amqp.ConnSASLPlain(clientUsername, clientPassword))
	}

	client, err := amqp.Dial(uri, opts...)
	if err != nil {
		return err
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal("Failed to close client:", err)
		}
	}()

	var ctx = context.Background()

	session, err := client.NewSession()
	if err != nil {
		return err
	}

	defer func() {
		if err := session.Close(ctx); err != nil {
			log.Fatal("Failed to close session:", err)
		}
	}()

	// create receivers

	telemetry, cleanup, err := createReceiver(session, tenant, "telemetry")
	if err != nil {
		return err
	}
	defer cleanup()
	event, cleanup, err := createReceiver(session, tenant, "event")
	if err != nil {
		return err
	}
	defer cleanup()

	// proceed

	fmt.Printf("Consumer running, press Ctrl+C to stop...")
	fmt.Println()

	// set up command reader

	var reader command.Reader
	if len(processCommands) > 0 {
		fmt.Printf("Enabling command reader (%s)...", commandReader)
		fmt.Println()
		reader, err = createCommandReader(processCommands)
		if err != nil {
			return err
		}
		if err := reader.Start(); err != nil {
			return err
		}
		defer func() {
			if err := reader.Stop(); err != nil {
			}
		}()
	} else {
		reader = nil
	}

	e := make(chan error)

	// run loops

	go func() { e <- runLoop(ctx, telemetry, session, reader, tenant) }()
	go func() { e <- runLoop(ctx, event, session, reader, tenant) }()

	// return error

	return <-e

}

func runLoop(ctx context.Context, receiver *amqp.Receiver, session *amqp.Session, reader command.Reader, tenant string) error {
	for {
		// Receive next message
		msg, err := receiver.Receive(ctx)
		if err != nil {
			return err
		}

		// Accept message
		if err := msg.Accept(); err != nil {
			return err
		}

		utils.PrintMessage(msg)
		if len(processCommands) > 0 {
			if err := processCommand(session, reader, tenant, msg); err != nil {
				log.Print("Failed to send command: ", err)
			}
		}
	}
}

func processCommand(session *amqp.Session, reader command.Reader, tenant string, msg *amqp.Message) error {
	ttd, ok := msg.ApplicationProperties["ttd"].(int32)

	if !ok {
		return nil
	}

	if ttd < 0 {
		return nil
	}

	deviceId, ok := msg.Annotations["device_id"].(string)
	if !ok || deviceId == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ttd)*time.Second)
	defer cancel()

	cmd := reader.Read(ctx, deviceId)

	if cmd == nil {
		fmt.Print("Timeout!")
		fmt.Println()
		return nil
	}

	// open sender

	sender, err := session.NewSender(
		amqp.LinkTargetAddress("command/" + tenant),
	)

	if err != nil {
		return err
	}

	// defer: close sender

	defer func() {
		if err := sender.Close(context.Background()); err != nil {
			log.Print("Failed to close sender: ", err)
		}
	}()

	// prepare payload
	var payload []byte
	if cmd.Payload == nil {
		payload = make([]byte, 0)
	} else {
		payload = cmd.Payload.Bytes()
	}

	// prepare message

	send := amqp.NewMessage(payload)
	send.Properties = &amqp.MessageProperties{
		Subject:     cmd.Command,
		ContentType: cmd.ContentType,
		To:          "command/" + tenant + "/" + deviceId,
	}

	// set message id

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	send.Properties.MessageID = amqp.UUID(id).String()

	// send message

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sender.Send(ctx, send); err != nil {
		return err
	}

	fmt.Println("Command delivered!")

	return nil
}
