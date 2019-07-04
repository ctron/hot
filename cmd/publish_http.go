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
	"bytes"
	"fmt"
	"net/http"
	neturl "net/url"
	"strconv"

	"github.com/ctron/hot/pkg/utils"
)

func publishHttp(messageType string, uri string, tenant string, deviceId string, authId string, password string, contentType string, payload string) error {

	url, err := neturl.Parse(uri)
	if err != nil {
		return err
	}

	url.Path = url.Path + neturl.PathEscape(messageType) + "/" + neturl.PathEscape(tenant) + "/" + deviceId
	fmt.Println("URL:", url)

	buf := bytes.NewBufferString(payload)

	tr := &http.Transport{
		TLSClientConfig: createTlsConfig(),
	}

	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("PUT", url.String(), buf)
	if err != nil {
		return err
	}

	request.SetBasicAuth(authId+"@"+tenant, password)

	if qos > 0 {
		request.Header.Set("QoS-Level", strconv.Itoa(int(qos)))
	}
	if ttd > 0 {
		request.Header.Set("hono-ttd", strconv.FormatUint(uint64(ttd), 10))
	}
	request.Header.Set("Content-Type", contentType)

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
