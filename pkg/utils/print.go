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

package utils

import (
	"fmt"

	"pack.ag/amqp"
)

func PrintTitle(title string) {
	fmt.Println("---------------------------------------------------------")
	fmt.Println("#", title)
	fmt.Println("---------------------------------------------------------")
}

func PrintEntry(k interface{}, v interface{}) {
	fmt.Printf("%s => %[2]v (%[2]T)", k, v)
	fmt.Println()
}

func PrintAnnotations(title string, data map[interface{}]interface{}) {
	if len(data) > 0 {
		PrintTitle(title)
		for k, v := range data {
			PrintEntry(k, v)
		}
	}
}

func PrintProperties(title string, data map[string]interface{}) {
	if len(data) > 0 {
		PrintTitle(title)
		for k, v := range data {
			PrintEntry(k, v)
		}
	}
}

func PrintMessageProperties(p *amqp.MessageProperties) {
	PrintTitle("Properties")

	fmt.Println("Content Encoding:", p.ContentEncoding)
	fmt.Println("Content Type:", p.ContentType)
	fmt.Println("Message ID:", p.MessageID)
	fmt.Println("Subject:", p.Subject)
	fmt.Println("To:", p.To)

}

func PrintStart() {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
}

func PrintEnd() {
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}

func PrintMessage(msg *amqp.Message) {

	PrintStart()

	PrintAnnotations("Annotations", msg.Annotations)
	PrintAnnotations("Delivery annotations", msg.DeliveryAnnotations)

	PrintMessageProperties(msg.Properties)
	PrintProperties("Application Properties", msg.ApplicationProperties)
	PrintAnnotations("Footer", msg.Footer)

	PrintTitle("Payload")

	fmt.Printf("%s", msg.GetData())
	fmt.Println()

	PrintEnd()

}
