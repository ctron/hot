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

package encoding

import (
	"bytes"
	"encoding/hex"
	"strings"
)

type PayloadEncoder interface {
	Encode(payloadString string) (*bytes.Buffer, error)
	GetMimeType() string
}

func CreateEncoder(contentType string) PayloadEncoder {

	switch strings.ToLower(contentType) {
	case "hex":
		return &HexPayloadEncoder{}
	default:
		return &RawPayloadEncoder{contentType}
	}

}

// Raw payload encoder, pass through string as blob

type RawPayloadEncoder struct {
	MimeType string
}

func (_ RawPayloadEncoder) Encode(payloadString string) (*bytes.Buffer, error) {
	return bytes.NewBufferString(payloadString), nil
}

func (e RawPayloadEncoder) GetMimeType() string {
	return e.MimeType
}

var _ PayloadEncoder = &RawPayloadEncoder{}

// Hex payload encoder

type HexPayloadEncoder struct {
}

func (_ HexPayloadEncoder) Encode(payloadString string) (*bytes.Buffer, error) {
	data, err := hex.DecodeString(payloadString)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

func (_ HexPayloadEncoder) GetMimeType() string {
	return "application/octet-stream"
}

var _ PayloadEncoder = &HexPayloadEncoder{}
