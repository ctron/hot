/*******************************************************************************
 * Copyright (c) 2020 Red Hat Inc
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
	"github.com/ctron/hot/pkg/encoding"
)

// static command reader

type StaticReader struct {
	commandName string
	contentType string
	payload     *bytes.Buffer
}

func NewStatic(commandName string, encoder encoding.PayloadEncoder, payload string) (*StaticReader, error) {

	encodedPayload, err := encoder.Encode(payload)
	if err != nil {
		return nil, err
	}

	return &StaticReader{
		commandName: commandName,
		contentType: encoder.GetMimeType(),
		payload:     encodedPayload,
	}, nil

}

func (_ *StaticReader) Start() error {
	return nil
}

func (_ *StaticReader) Stop() error {
	return nil
}

func (s *StaticReader) Read(_ context.Context, _ string) *Command {
	return &Command{
		Command:     s.commandName,
		ContentType: s.contentType,
		Payload:     s.payload,
	}
}

var _ Reader = &StaticReader{}
