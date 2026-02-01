package ndjson

import (
	"errors"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{name: "nil value", input: nil},
		{name: "string", input: "hello"},
		{name: "int", input: 42},
		{name: "struct", input: struct{ Name string }{Name: "test"}},
		{name: "map", input: map[string]int{"a": 1, "b": 2}},
		{name: "slice", input: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)

			if data != nil {
				t.Errorf("Marshal(%v) returned data = %v, want nil", tt.input, data)
			}

			if !errors.Is(err, ErrNotImplemented) {
				t.Errorf("Marshal(%v) returned err = %v, want ErrNotImplemented", tt.input, err)
			}
		})
	}
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
