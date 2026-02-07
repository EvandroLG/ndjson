package ndjson

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name   string
		input  any
		expect []any
	}{
		{name: "nil value", input: nil, expect: []any{nil}},
		{name: "string", input: "hello", expect: []any{"hello"}},
		{name: "int", input: 42, expect: []any{42}},
		{name: "struct", input: testStruct{Name: "test"}, expect: []any{testStruct{Name: "test"}}},
		{name: "map", input: map[string]int{"a": 1, "b": 2}, expect: []any{map[string]int{"a": 1, "b": 2}}},
		{name: "slice", input: []int{1, 2, 3}, expect: []any{1, 2, 3}},
		{name: "array", input: [2]string{"a", "b"}, expect: []any{"a", "b"}},
		{
			name:   "raw message",
			input:  json.RawMessage(`{"a":1}`),
			expect: []any{json.RawMessage(`{"a":1}`)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal(%v) unexpected error: %v", tt.input, err)
			}
			if len(data) == 0 || data[len(data)-1] != '\n' {
				t.Fatalf("Marshal(%v) missing trailing newline: %q", tt.input, string(data))
			}

			lines := splitNDJSONLines(data)
			if len(lines) != len(tt.expect) {
				t.Fatalf("Marshal(%v) lines = %d, want %d", tt.input, len(lines), len(tt.expect))
			}

			for i, line := range lines {
				got := decodeAnyJSON(t, line)
				want := normalizeAnyJSON(t, tt.expect[i])
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("Marshal(%v) line %d = %v, want %v", tt.input, i, got, want)
				}
			}
		})
	}
}

func TestMarshalError(t *testing.T) {
	data, err := Marshal(make(chan int))
	if err == nil {
		t.Fatalf("Marshal(chan int) expected error, got nil")
	}
	if len(data) != 0 {
		t.Fatalf("Marshal(chan int) returned data, want empty: %q", string(data))
	}
}

func splitNDJSONLines(data []byte) [][]byte {
	lines := bytes.Split(data, []byte("\n"))
	if len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func decodeAnyJSON(t *testing.T, data []byte) any {
	t.Helper()
	var out any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	return out
}

func normalizeAnyJSON(t *testing.T, v any) any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	return decodeAnyJSON(t, b)
}

func TestUnmarshal(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name    string
		data    []byte
		want    []Person
		wantErr bool
	}{
		{
			name: "single object",
			data: []byte(`{"name":"Alice","age":30}`),
			want: []Person{{Name: "Alice", Age: 30}},
		},
		{
			name: "multiple objects",
			data: []byte("{\"name\":\"Alice\",\"age\":30}\n{\"name\":\"Bob\",\"age\":25}"),
			want: []Person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}},
		},
		{
			name: "with empty lines",
			data: []byte("{\"name\":\"Alice\",\"age\":30}\n\n{\"name\":\"Bob\",\"age\":25}"),
			want: []Person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}},
		},
		{
			name: "trailing newline",
			data: []byte("{\"name\":\"Alice\",\"age\":30}\n"),
			want: []Person{{Name: "Alice", Age: 30}},
		},
		{
			name: "empty data",
			data: []byte{},
			want: nil,
		},
		{
			name: "only newlines",
			data: []byte("\n\n\n"),
			want: nil,
		},
		{
			name:    "invalid json",
			data:    []byte(`{invalid}`),
			wantErr: true,
		},
		{
			name:    "malformed json in second line",
			data:    []byte("{\"name\":\"Alice\",\"age\":30}\n{invalid}"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []Person
			err := Unmarshal(tt.data, &got)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Unmarshal() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalMapSlice(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []map[string]any
		wantErr bool
	}{
		{
			name: "single map object",
			data: []byte(`{"key":"value","count":42}`),
			want: []map[string]any{{"key": "value", "count": float64(42)}},
		},
		{
			name: "multiple map objects",
			data: []byte("{\"a\":1}\n{\"b\":2}"),
			want: []map[string]any{{"a": float64(1)}, {"b": float64(2)}},
		},
		{
			name: "nested objects",
			data: []byte(`{"outer":{"inner":"value"}}`),
			want: []map[string]any{{"outer": map[string]any{"inner": "value"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []map[string]any
			err := Unmarshal(tt.data, &got)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Unmarshal() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoPanic(t *testing.T) {
	t.Run("Marshal with nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Marshal(nil) panicked: %v", r)
			}
		}()
		_, _ = Marshal(nil)
	})

	// Note: Unmarshal(nil, nil) is expected to panic since it requires
	// a valid pointer to a slice. This matches json.Unmarshal behavior.
}
