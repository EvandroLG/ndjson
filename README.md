# ndjson

A Go package for encoding and decoding Newline Delimited JSON (NDJSON).

> **Note:** This package is currently a stub. `Marshal` and `Unmarshal` return `ErrNotImplemented`.

## Installation

```bash
go get github.com/evandrolg/ndjson
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/evandrolg/ndjson"
)

func main() {
	// Marshal (not yet implemented)
	data, err := ndjson.Marshal(map[string]string{"hello": "world"})
	if err != nil {
		log.Printf("Marshal error: %v", err) // ndjson: not implemented
	}

	// Unmarshal (not yet implemented)
	var results []map[string]any
	err = ndjson.Unmarshal([]byte(`{"a":1}\n{"b":2}`), &results)
	if err != nil {
		log.Printf("Unmarshal error: %v", err) // ndjson: not implemented
	}

	fmt.Println(data) // nil
}
```

## License

MIT
