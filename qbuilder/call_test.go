package qbuilder

import (
	"testing"
)

func TestCallBuilder_ToSQL(t *testing.T) {
	tests := []struct {
		name     string
		builder  func() *CallBuilder
		expected string
	}{
		{
			name: "No arguments",
			builder: func() *CallBuilder {
				return new(CallBuilder).Call("my_function")
			},
			expected: "CALL my_function",
		},
		{
			name: "One argument",
			builder: func() *CallBuilder {
				return new(CallBuilder).Call("my_function").Argument("$1")
			},
			expected: "CALL my_function ($1)",
		},
		{
			name: "Multiple arguments",
			builder: func() *CallBuilder {
				return new(CallBuilder).Call("my_function").Argument("$1").Argument("$2")
			},
			expected: "CALL my_function ($1, $2)",
		},
		{
			name: "Arguments using Arguments method",
			builder: func() *CallBuilder {
				return new(CallBuilder).Call("my_function").Arguments("$1", "$2", "$3")
			},
			expected: "CALL my_function ($1, $2, $3)",
		},
		{
			name: "Arguments with different types",
			builder: func() *CallBuilder {
				return new(CallBuilder).Call("my_function").Argument("$1").Argument("$2").Argument("$3")
			},
			expected: "CALL my_function ($1, $2, $3)",
		},
		{
			name: "Empty function call with no arguments",
			builder: func() *CallBuilder {
				return new(CallBuilder).Call("empty_function")
			},
			expected: "CALL empty_function",
		},
		{
			name: "Clear arguments",
			builder: func() *CallBuilder {
				b := new(CallBuilder).Call("my_function").Argument("$1").Argument("$2")
				b.Clear()
				return b
			},
			expected: "CALL my_function",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.builder()
			result := b.ToSQL()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
