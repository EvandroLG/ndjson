# ndjson &middot;[![CI](https://github.com/evandrolg/ndjson/actions/workflows/ci.yml/badge.svg)](https://github.com/evandrolg/ndjson/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/evandrolg/ndjson.svg)](https://pkg.go.dev/github.com/evandrolg/ndjson) [![Go Report Card](https://goreportcard.com/badge/github.com/evandrolg/ndjson)](https://goreportcard.com/report/github.com/evandrolg/ndjson) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)

A Go package for encoding and decoding Newline Delimited JSON (NDJSON).

Marshal encodes values to NDJSON and Unmarshal decodes NDJSON into a slice.

## Installation

```bash
go get github.com/evandrolg/ndjson
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/evandrolg/ndjson"
)

func main() {
	data, err := ndjson.Marshal([]map[string]any{
		{"name": "Alice"},
		{"name": "Bob"},
	})
	if err != nil {
		fmt.Errorf("Marshal error: %v\n", err)
		return
	}

	var results []map[string]any
	err = ndjson.Unmarshal(data, &results)
	if err != nil {
		fmt.Errorf("Unmarshal error: %v\n", err)
		return
	}

	fmt.Println(results)
}
```

## License

[MIT](./LICENSE)
