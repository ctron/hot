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

package async

import (
	"context"
	"io"
)

type Reader interface {
	Close() error

	ReadLine(ctx context.Context) *string
}

type channelReader struct {
	reader io.Reader
	data   chan string
}

var _ Reader = &channelReader{}

func NewReader(reader io.Reader) Reader {
	data := make(chan string, 1)

	ChannelReader(reader, data)

	return &channelReader{
		reader: reader,
		data:   data,
	}
}

func (r *channelReader) Close() error {
	close(r.data)
	return nil
}

func (r *channelReader) ReadLine(ctx context.Context) *string {
	select {
	case s := <-r.data:
		return &s
	case <-ctx.Done():
		return nil
	}
}
