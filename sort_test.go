package rel

import (
	"reflect"
	"testing"
)

func TestParseSort(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		allowed  map[string]struct{}
		expected Sort
	}{
		{
			name:     "Empty query",
			query:    "",
			allowed:  nil,
			expected: nil,
		},
		{
			name:     "Single column, default order",
			query:    "name",
			allowed:  nil,
			expected: Sort{{Column: "name", Order: OrderAsc}},
		},
		{
			name:     "Single column, descending order",
			query:    "created_at:desc",
			allowed:  nil,
			expected: Sort{{Column: "created_at", Order: OrderDesc}},
		},
		{
			name:     "Multiple columns with mixed order",
			query:    "name,created_at:desc",
			allowed:  nil,
			expected: Sort{{Column: "name", Order: OrderAsc}, {Column: "created_at", Order: OrderDesc}},
		},
		{
			name:     "Ignore disallowed columns",
			query:    "name,secret,created_at:desc",
			allowed:  map[string]struct{}{"name": {}, "created_at": {}},
			expected: Sort{{Column: "name", Order: OrderAsc}, {Column: "created_at", Order: OrderDesc}},
		},
		{
			name:     "Column with empty key",
			query:    ":desc",
			allowed:  nil,
			expected: nil,
		},
		{
			name:     "Column with spaces",
			query:    "  name  ,created_at: desc  ",
			allowed:  nil,
			expected: Sort{{Column: "name", Order: OrderAsc}, {Column: "created_at", Order: OrderDesc}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseSort(tt.query, tt.allowed)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %+v, expected %+v", result, tt.expected)
			}
		})
	}
}
