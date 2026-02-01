package ndjson

import (
	"errors"
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
	tests := []struct {
		name   string
		data   []byte
		target any
	}{
		{name: "nil data and target", data: nil, target: nil},
		{name: "empty data", data: []byte{}, target: new([]map[string]any)},
		{name: "valid ndjson", data: []byte(`{"a":1}\n{"b":2}`), target: new([]map[string]any)},
		{name: "into struct slice", data: []byte(`{"name":"test"}`), target: new([]struct{ Name string })},
		{name: "into map", data: []byte(`{"key":"value"}`), target: new(map[string]string)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal(tt.data, tt.target)

			if !errors.Is(err, ErrNotImplemented) {
				t.Errorf("Unmarshal(%v, %v) returned err = %v, want ErrNotImplemented", tt.data, tt.target, err)
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

	t.Run("Unmarshal with nil data and target", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Unmarshal(nil, nil) panicked: %v", r)
			}
		}()
		_ = Unmarshal(nil, nil)
	})
}
