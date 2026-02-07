package ndjson

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// Marshal encodes v as one or more NDJSON lines.
//
// The input v can be any value serializable to JSON.
// If v is a slice or array, each element is encoded as a separate JSON value with a trailing newline.
// If v is any other value, it is encoded as a single JSON value with a trailing newline.
//
// The returned byte slice contains the NDJSON data.
// Marshal returns an error if any value cannot be encoded.
func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	if v == nil {
		if err := enc.Encode(nil); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	rv := reflect.ValueOf(v)
	rt := rv.Type()

	// json.RawMessage (and any Marshaler) should be encoded as a single value
	if _, ok := v.(json.Marshaler); ok {
		if err := enc.Encode(v); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	// byte slices/arrays (named or unnamed) should be a single value
	if (rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) && rt.Elem().Kind() == reflect.Uint8 {
		if err := enc.Encode(v); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i).Interface()
			if err := enc.Encode(elem); err != nil {
				return nil, err
			}
		}
	default:
		if err := enc.Encode(v); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// Unmarshal decodes NDJSON data into v.
//
// The data parameter should contain one or more newline-delimited JSON objects.
// The v parameter must be a pointer to a slice or other suitable container
// that can hold the decoded values.
func Unmarshal(data []byte, v any) error {
	fragments := bytes.Split(data, []byte("\n"))

	sliceValue := reflect.ValueOf(v).Elem()
	elemType := sliceValue.Type().Elem()

	for _, fragment := range fragments {
		if len(fragment) == 0 {
			continue
		}

		elem := reflect.New(elemType).Interface()
		err := json.Unmarshal(fragment, elem)
		if err != nil {
			return err
		}

		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(elem).Elem()))
	}
	return nil
}
