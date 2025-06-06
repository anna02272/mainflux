// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"encoding/json"
	"strings"

	"github.com/MainfluxLabs/mainflux/pkg/errors"
	protomfx "github.com/MainfluxLabs/mainflux/pkg/proto"
)

var (
	// ErrTransform represents an error during parsing message.
	ErrTransform = errors.New("unable to parse JSON object")
	// ErrInvalidTimeField represents the use an invalid time field.
	ErrInvalidTimeField = errors.New("invalid time field")

	errInvalidFormat     = errors.New("invalid JSON object")
	errInvalidNestedJSON = errors.New("invalid nested JSON object")
)

func TransformPayload(transformer protomfx.Transformer, msg *protomfx.Message) error {
	var payload interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return errors.Wrap(ErrTransform, err)
	}

	extractedPayload := extractPayload(payload, transformer.DataField)

	switch p := extractedPayload.(type) {
	case map[string]interface{}:
		ts, err := transformTimeField(p, transformer)
		if err != nil {
			return errors.Wrap(ErrInvalidTimeField, err)
		}

		if ts != 0 {
			msg.Created = ts
		}

		formattedPayload := filterPayloadFields(p, transformer.DataFilters)
		data, err := json.Marshal(formattedPayload)
		if err != nil {
			return err
		}
		msg.Payload = data
		return nil

	case []interface{}:
		var payloads []map[string]interface{}
		for _, val := range p {
			v, ok := val.(map[string]interface{})
			if !ok {
				return errors.Wrap(ErrTransform, errInvalidNestedJSON)
			}

			ts, err := transformTimeField(v, transformer)
			if err != nil {
				return errors.Wrap(ErrInvalidTimeField, err)
			}

			if ts != 0 {
				v["Created"] = ts
			}

			formattedPayload := filterPayloadFields(v, transformer.DataFilters)
			payloads = append(payloads, formattedPayload)
		}

		data, err := json.Marshal(payloads)
		if err != nil {
			return err
		}
		msg.Payload = data
		return nil

	default:
		return errors.Wrap(ErrTransform, errInvalidFormat)
	}
}

func transformTimeField(payload interface{}, transformer protomfx.Transformer) (int64, error) {
	if transformer.TimeField == "" {
		return 0, nil
	}

	val := payload
	keys := strings.Split(transformer.TimeField, ".")

	for _, k := range keys {
		switch v := val.(type) {
		case map[string]interface{}:
			current, exists := v[k]
			if !exists {
				return 0, nil
			}
			val = current

		case []interface{}:
			if len(v) == 0 {
				return 0, nil
			}
			firstItem, ok := v[0].(map[string]interface{})
			if !ok {
				return 0, nil
			}
			current, exists := firstItem[k]
			if !exists {
				return 0, nil
			}
			val = current

		default:
			return 0, nil
		}
	}

	t, err := parseTimestamp(transformer.TimeFormat, val, transformer.TimeLocation)
	if err != nil {
		return 0, err
	}
	return t.UnixNano(), nil
}

func extractPayload(payload interface{}, dataField string) interface{} {
	if dataField == "" {
		return payload
	}

	value := payload
	keys := strings.Split(dataField, ".")

	for _, k := range keys {
		switch v := value.(type) {
		case map[string]interface{}:
			current, exists := v[k]
			if !exists {
				return nil
			}
			value = current

		case []interface{}:
			if len(v) == 0 {
				return nil
			}
			firstItem, ok := v[0].(map[string]interface{})
			if !ok {
				return nil
			}
			current, exists := firstItem[k]
			if !exists {
				return nil
			}
			value = current

		default:
			return nil
		}
	}

	return value
}

func filterPayloadFields(payload map[string]interface{}, dataFilters []string) map[string]interface{} {
	if len(dataFilters) == 0 {
		return payload
	}

	filteredPayload := make(map[string]interface{})

	for _, key := range dataFilters {
		// Split nested path
		keys := strings.Split(key, ".")
		var value interface{} = payload

		// Traverse nested structure
		for _, k := range keys {
			current, ok := value.(map[string]interface{})
			if !ok {
				value = nil
				break
			}

			v, exists := current[k]
			if !exists {
				value = nil
				break
			}
			value = v
		}

		if value != nil {
			filteredPayload[key] = value
		}
	}

	return filteredPayload
}
