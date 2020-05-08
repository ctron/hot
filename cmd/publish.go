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

type CommonPublishInformation struct {
	MessageType string
	URI         string
	Tenant      string
	DeviceId    string

	AuthenticationId string
	Username         string
	Password         string
}

func (c CommonPublishInformation) HasUsernamePassword() bool {
	return c.Password != "" || c.Username != "" || c.AuthenticationId != ""
}

func (c CommonPublishInformation) EffectiveUsername() string {
	if c.Username != "" {
		return c.Username
	} else {
		return c.AuthenticationId + "@" + c.Tenant
	}
}
