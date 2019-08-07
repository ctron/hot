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
	"github.com/ctron/hot/pkg/encoding"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttPublishInformation struct {
	CommonPublishInformation

	QoS uint8
}

func publishMqtt(info MqttPublishInformation, encoder encoding.PayloadEncoder, payload string) error {

	opts := MQTT.NewClientOptions()
	opts.AddBroker(info.URI)
	opts.SetClientID(info.DeviceId)
	opts.SetUsername(info.Username())
	opts.SetPassword(info.Password)
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

	// publish

	token := client.Publish(topic, byte(info.QoS), false, buf.Bytes())
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	// done

	return nil
}
