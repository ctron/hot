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
	"fmt"
	"github.com/ctron/hot/pkg/encoding"
	"github.com/ctron/hot/pkg/utils"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type MqttPublishInformation struct {
	CommonPublishInformation

	QoS uint8
}

func publishMqtt(info MqttPublishInformation, encoder encoding.PayloadEncoder, payload string) error {

	if err := validateMessageType(info.MessageType); err != nil {
		return err
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(info.URI)
	opts.SetClientID(info.DeviceId)
	if info.HasUsernamePassword() {
		opts.SetUsername(info.EffectiveUsername())
		opts.SetPassword(info.Password)
	}
	opts.SetCleanSession(true)
	opts.SetTLSConfig(createTlsConfig())

	topic := info.MessageType

	buf, err := encoder.Encode(payload)
	if err != nil {
		return err
	}

	// connect to MQTT endpoint

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	// defer - close connection

	defer client.Disconnect(250)

	// set up command handler

	var response chan struct{} = nil
	if ttd > 0 {
		fmt.Println("Subscribing to command topic ... ")
		response = make(chan struct{})
		subToken := client.Subscribe("command///req/#", 0, func(client MQTT.Client, message MQTT.Message) {
			utils.PrintStart()

			utils.PrintTitle("Headers")
			utils.PrintEntry("Topic", message.Topic())
			utils.PrintEntry("QoS", message.Qos())
			utils.PrintEntry("MessageID", message.MessageID())
			utils.PrintEntry("Ack", message.Ack)
			utils.PrintEntry("Duplicate", message.Duplicate())
			utils.PrintEntry("Retained", message.Retained())

			utils.PrintTitle("Payload")
			fmt.Println(message.Payload())

			utils.PrintEnd()

			response <- struct{}{}
		})
		if subToken.WaitTimeout(time.Second*10) && subToken.Error() != nil {
			return subToken.Error()
		}
		fmt.Println("Subscribing to command topic ... OK!")
	}

	// publish

	token := client.Publish(topic, byte(info.QoS), false, buf.Bytes())
	if token.WaitTimeout(time.Second*10) && token.Error() != nil {
		return token.Error()
	}

	if response != nil {
		select {
		case <-time.After(time.Second * time.Duration(ttd)):
			fmt.Println("No command received in time")
		case <-response:
		}
	}

	// done

	return nil
}
