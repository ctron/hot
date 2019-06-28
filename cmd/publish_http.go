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

	fmt.Println(body.String())

	if err := response.Body.Close(); err != nil {
		fmt.Printf("Failed to close response: %v", err)
	}

	return nil
}
