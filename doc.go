// Package ndjson provides encoding and decoding of Newline Delimited JSON (NDJSON).
//
// NDJSON is a convenient format for storing or streaming structured data
// that may be processed one record at a time. Each line in an NDJSON file
// is a valid JSON value, typically an object. Lines are separated by '\n'.
//
// Example NDJSON content:
//
//	{"name":"Alice","age":30}
//	{"name":"Bob","age":25}
//	{"name":"Charlie","age":35}
//
// This package provides Marshal and Unmarshal functions similar to encoding/json,
// but designed for the NDJSON format.
//
// Marshal encodes values as NDJSON:
//   - Slices and arrays are encoded one element per line.
//   - Other values are encoded as a single JSON value with a trailing newline.
//   - Values implementing json.Marshaler and byte slices/arrays are encoded as a
//     single JSON value (not split into lines).
//
// Unmarshal decodes NDJSON into a destination. The destination must be a
// pointer to a slice (values are appended). Empty lines are ignored, and
// JSON decode errors include the line number.
package ndjson
