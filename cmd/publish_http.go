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
	"bytes"
	"fmt"
	"net/http"
	neturl "net/url"
	"strconv"

	"github.com/ctron/hot/pkg/encoding"
	"github.com/ctron/hot/pkg/utils"
)

type HttpPublishInformation struct {
	CommonPublishInformation

	QoS uint8
}

func publishHttp(info HttpPublishInformation, encoder encoding.PayloadEncoder, payload string) error {

	if err := validateMessageType(info.MessageType); err != nil {
		return err
	}

	url, err := neturl.Parse(info.URI)
	if err != nil {
		return err
	}

	url.Path = url.Path + neturl.PathEscape(info.MessageType) + "/" + neturl.PathEscape(info.Tenant) + "/" + info.DeviceId
	fmt.Println("URL:", url)

	buf, err := encoder.Encode(payload)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: createTlsConfig(),
	}

	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("PUT", url.String(), buf)
	if err != nil {
		return err
	}

	if info.HasUsernamePassword() {
		request.SetBasicAuth(info.EffectiveUsername(), info.Password)
	}

	if qos > 0 {
		request.Header.Set("QoS-Level", strconv.Itoa(int(qos)))
	}
	if ttd > 0 {
		request.Header.Set("hono-ttd", strconv.FormatUint(uint64(ttd), 10))
	}
	request.Header.Set("Content-Type", encoder.GetMimeType())

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	fmt.Printf("Publish result: %s", response.Status)
	fmt.Println()

	body := new(bytes.Buffer)
	if _, err := body.ReadFrom(response.Body); err != nil {
		return err
	}

	utils.PrintStart()

	utils.PrintTitle("Headers")
	for k, v := range response.Header {
		if len(v) == 1 {
			utils.PrintEntry(k, v[0])
		} else {
			utils.PrintEntry(k, v)
		}
	}

	if body.Len() > 0 {
		utils.PrintTitle("Payload")
		fmt.Println(body.String())
	}
	utils.PrintEnd()

	if err := response.Body.Close(); err != nil {
		fmt.Printf("Failed to close response: %v", err)
	}

	return nil
}
