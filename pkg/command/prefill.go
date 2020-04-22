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
	"context"
	"fmt"
	"io"

	"github.com/ctron/hot/pkg/async"

	"github.com/ctron/hot/pkg/encoding"
)

// pre fill reader

type PreFillReader struct {
	CommandName string
	Stream      io.Reader
	Prompt      io.Writer
	Encoder     encoding.PayloadEncoder
	cmd         chan Command
}

func (r *PreFillReader) prompt() {
	if r.Prompt != nil {
		_, _ = fmt.Fprint(r.Prompt, "Next command: ")
	}
}

func (r *PreFillReader) Start() error {

	r.cmd = make(chan Command)

	r.prompt()

	async.CallbackReader(r.Stream, func(s *string, e error) bool {
		if e != nil {
			return false
		}

		defer r.prompt()

		if s == nil || *s == "" {
			return true
		}

		payload, err := r.Encoder.Encode(*s)
		if err != nil {
			if r.Prompt != nil {
				_, _ = fmt.Fprintf(r.Prompt, "Failed to convert payload: %v", err)
				_, _ = fmt.Fprintln(r.Prompt)
			}
			// payload format error
			return true
		}

		cmd := Command{
			Command:     r.CommandName,
			ContentType: r.Encoder.GetMimeType(),
			Payload:     payload,
		}

		if r.Prompt != nil {
			_, _ = fmt.Fprintf(r.Prompt, "Next command set to: %+v", cmd)
			_, _ = fmt.Fprintln(r.Prompt)
		}

		r.cmd <- cmd

		if r.Prompt != nil {
			_, _ = fmt.Fprintln(r.Prompt, "Command released")
		}

		return true
	})

	return nil
}

func (r *PreFillReader) Stop() error {
	return nil
}

func (r *PreFillReader) Read(ctx context.Context, _ string) *Command {
	select {
	case cmd := <-r.cmd:
		return &cmd
	case <-ctx.Done():
		return nil
	}
}

var _ Reader = &PreFillReader{}
