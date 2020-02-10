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

func createCommandReader() command.Reader {
	switch strings.ToLower(commandReader) {
	case "prefill":
		return &command.PreFillReader{
			Prompt:  os.Stdout,
			Stream:  os.Stdin,
			Encoder: getEncoder(),
		}
	case "ondemand":
		return &command.OnDemandReader{
			Prompt:  os.Stdout,
			Stream:  os.Stdin,
			Encoder: getEncoder(),
		}
	default:
		panic(fmt.Errorf("unknown command reader: %s", commandReader))
	}
}

func consume(messageType string, uri string, tenant string) error {

	fmt.Printf("Consuming %s from %s ...", messageType, uri)
	fmt.Println()

	opts := make([]amqp.ConnOption, 0)

	//Enable TLS if required
	if tlsConfig != 0 {
		opts = append(opts, amqp.ConnTLSConfig(createTlsConfig()))
	}
	
	//Enable Client credentials if avaliable
	if(clientUsername != "" && clientPassword !=""){
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

	receiver, err := session.NewReceiver(
		amqp.LinkSourceAddress(fmt.Sprintf("%s/%s", messageType, tenant)),
		amqp.LinkCredit(10),
	)
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		if err := receiver.Close(ctx); err != nil {
			log.Fatal("Failed to close receiver: ", err)
		}
		cancel()
	}()

	fmt.Printf("Consumer running, press Ctrl+C to stop...")
	fmt.Println()

	// set up command reader

	reader := createCommandReader()
	if err := reader.Start(); err != nil {
		return err
	}
	defer func() {
		if err := reader.Stop(); err != nil {
		}
	}()

	// run loop

	for {
		// Receive next message
		msg, err := receiver.Receive(ctx)
		if err != nil {
			return err
		}

		// Accept message
		if err := msg.Accept(); err != nil {
			return nil
		}

		utils.PrintMessage(msg)
		if processCommands {
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
		amqp.LinkTargetAddress("control/" + tenant + "/" + deviceId),
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
		To:          "control/" + tenant + "/" + deviceId,
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
