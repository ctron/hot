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

package main

import (
	"fmt"
)

func validateMessageType(messageType string) error {
	switch messageType {
	case "telemetry": // ok
	case "event": // ok
	default:
		return fmt.Errorf("message type '%s' is not supported", messageType)
	}

	return nil
}
