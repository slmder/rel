package qbuilder

import "testing"

func TestUpdateBuilder_ToSQL(t *testing.T) {
	tests := []struct {
		name     string
		builder  *UpdateBuilder
		expected string
	}{
		{
			name: "Basic Update",
			builder: Update("users").
				Set("name", "'John'"),
			expected: "UPDATE users SET name = 'John'",
		},
		{
			name: "Update With Where",
			builder: Update("users").
				Set("name", "'John'").
				Where("id = 1"),
			expected: "UPDATE users SET name = 'John' WHERE id = 1",
		},
		{
			name: "Update With Multiple Set",
			builder: Update("users").
				Set("name", "'John'").
				Set("age", "30"),
			expected: "UPDATE users SET name = 'John', age = 30",
		},
		{
			name: "Update With Multiple Where",
			builder: Update("users").
				Set("name", "'John'").
				Where("id = 1").
				AndWhere("age > 18"),
			expected: "UPDATE users SET name = 'John' WHERE id = 1 AND age > 18",
		},
		{
			name: "Update With Alias",
			builder: Update("users").
				As("u").
				Set("u.name", "'John'"),
			expected: "UPDATE users AS u SET u.name = 'John'",
		},
		{
			name: "Update With Returning",
			builder: Update("users").
				Set("name", "'John'").
				Returning("id", "name"),
			expected: "UPDATE users SET name = 'John' RETURNING id, name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.builder.ToSQL()
			if got != tt.expected {
				t.Errorf("expected: %s, got: %s", tt.expected, got)
			}
		})
	}
}
