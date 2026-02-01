package ndjson

import "errors"

// ErrNotImplemented is returned by Marshal and Unmarshal until
// the package implementation is complete.
var ErrNotImplemented = errors.New("ndjson: not implemented")
