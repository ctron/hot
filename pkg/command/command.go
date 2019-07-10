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

package command

import (
	"bytes"
	"context"
)

type Command struct {
	Command     string
	ContentType string
	Payload     *bytes.Buffer
}

type Reader interface {
	Start() error
	Stop() error

	Read(ctx context.Context, deviceId string) *Command
}

func GetReader() Reader {
	return &OnDemandReader{}
}
