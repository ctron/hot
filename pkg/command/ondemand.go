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
	"math"
	"time"

	"github.com/ctron/hot/pkg/async"

	"github.com/ctron/hot/pkg/encoding"
)

// on demand reader

type OnDemandReader struct {
	CommandName string
	Stream      io.Reader
	Prompt      io.Writer
	Encoder     encoding.PayloadEncoder
	reader      async.Reader
}

func (r *OnDemandReader) Start() error {
	r.reader = async.NewReader(r.Stream)
	return nil
}

func (r *OnDemandReader) Stop() error {
	return r.reader.Close()
}

func (r *OnDemandReader) Read(ctx context.Context, _ string) *Command {

	if r.Prompt != nil {
		deadline, ok := ctx.Deadline()
		if !ok {
			_, _ = fmt.Fprint(r.Prompt, "Command: ")
		} else {
			var rem float64 = float64(deadline.Sub(time.Now())) / float64(time.Second)
			_, _ = fmt.Fprintf(r.Prompt, "Command (%.0fs): ", math.Ceil(rem))
		}
	}

	payload := r.reader.ReadLine(ctx)

	if payload == nil {
		return nil
	} else {

		p, err := r.Encoder.Encode(*payload)
		if err != nil {
			return nil
		}

		return &Command{
			Command:     r.CommandName,
			ContentType: r.Encoder.GetMimeType(),
			Payload:     p,
		}
	}
}

var _ Reader = &OnDemandReader{}
