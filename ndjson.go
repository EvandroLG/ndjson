package ndjson

// Marshal encodes v as a single NDJSON line.
//
// Each call to Marshal produces one JSON object followed by a newline.
// The input v can be any value that is serializable to JSON.
//
// Note: This function is not yet implemented and will return ErrNotImplemented.
func Marshal(v any) ([]byte, error) {
	return nil, ErrNotImplemented
}

// Unmarshal decodes NDJSON data into v.
//
// The data parameter should contain one or more newline-delimited JSON objects.
// The v parameter must be a pointer to a slice or other suitable container
// that can hold the decoded values.
//
// Note: This function is not yet implemented and will return ErrNotImplemented.
func Unmarshal(data []byte, v any) error {
	return ErrNotImplemented
}
