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

import "context"

// noop command reader

type NoopReader struct {
}

func (_ *NoopReader) Start() error {
	return nil
}

func (_ *NoopReader) Stop() error {
	return nil
}

func (_ *NoopReader) Read(ctx context.Context, deviceId string) *Command {
	return nil
}

var _ Reader = &NoopReader{}
