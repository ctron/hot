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
	"bufio"
	"os"
	"time"
)

type CommandReader interface {
	ReadCommand(timeout time.Duration) *string
}

// noop command reader

type NoopCommandReader struct {
}

func (_ *NoopCommandReader) ReadCommand(timeout time.Duration) *string {
	return nil
}

var _ CommandReader = &NoopCommandReader{}

// stdin command reader

type StdinCommandReader struct {
}

func (_ *StdinCommandReader) ReadCommand(timeout time.Duration) *string {

	s := make(chan string)
	e := make(chan error)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			e <- err
		} else {
			s <- line
		}
		close(s)
		close(e)
	}()

	select {
	case line := <-s:
		return &line
	case _ = <-e:
		return nil
	case <-time.After(timeout):
		return nil
	}

}

var _ CommandReader = &StdinCommandReader{}
