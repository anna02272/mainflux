// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package senml

import (
	"encoding/json"

	"github.com/MainfluxLabs/mainflux/pkg/errors"
	protomfx "github.com/MainfluxLabs/mainflux/pkg/proto"
	"github.com/MainfluxLabs/senml"
)

const (
	// JSON represents SenML in JSON format content type.
	JSON = "application/senml+json"
	// CBOR represents SenML in CBOR format content type.
	CBOR = "application/senml+cbor"
)

var (
	errDecode    = errors.New("failed to decode senml")
	errNormalize = errors.New("failed to normalize senml")
)

var formats = map[string]senml.Format{
	JSON: senml.JSON,
	CBOR: senml.CBOR,
}

func TransformPayload(msg *protomfx.Message) error {
	contentFormat := msg.ContentType
	format, ok := formats[contentFormat]
	if !ok {
		format = formats[JSON]
	}

	raw, err := senml.Decode(msg.Payload, format)
	if err != nil {
		return errors.Wrap(errDecode, err)
	}

	normalized, err := senml.Normalize(raw)
	if err != nil {
		return errors.Wrap(errNormalize, err)
	}

	payloads := make([]Message, len(normalized.Records))
	for i, v := range normalized.Records {
		// Use reception timestamp if SenML message Time is missing
		t := v.Time
		if t == 0 {
			// Convert the Unix timestamp in nanoseconds to float64
			t = float64(msg.Created) / 1e9
		}

		payloads[i] = Message{
			Name:        v.Name,
			Unit:        v.Unit,
			Time:        t,
			UpdateTime:  v.UpdateTime,
			Value:       v.Value,
			BoolValue:   v.BoolValue,
			DataValue:   v.DataValue,
			StringValue: v.StringValue,
			Sum:         v.Sum,
		}
	}

	data, err := json.Marshal(payloads)
	if err != nil {
		return err
	}
	msg.Payload = data

	return nil
}
