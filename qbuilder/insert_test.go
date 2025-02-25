package qbuilder

import (
	"testing"
)

func TestInsertBuilder_ToSQL(t *testing.T) {
	tests := []struct {
		name     string
		builder  *InsertBuilder
		expected string
	}{
		{
			name:     "Simple INSERT",
			builder:  Insert("users").Columns("id", "name").Values([]string{"1", "'John'"}),
			expected: "INSERT INTO users (id, name) VALUES (1, 'John')",
		},
		{
			name:     "INSERT with multiple values",
			builder:  Insert("users").Columns("id", "name").Values([]string{"1", "'John'"}, []string{"2", "'Jane'"}),
			expected: "INSERT INTO users (id, name) VALUES (1, 'John'), (2, 'Jane')",
		},
		{
			name:     "INSERT with alias",
			builder:  Insert("users").As("u").Columns("id", "name").Values([]string{"1", "'John'"}),
			expected: "INSERT INTO users AS u (id, name) VALUES (1, 'John')",
		},
		{
			name:     "INSERT with ON CONFLICT DO NOTHING",
			builder:  Insert("users").Columns("id", "name").Values([]string{"1", "'John'"}).OnConflict("id", false).DoNothing(),
			expected: "INSERT INTO users (id, name) VALUES (1, 'John') ON CONFLICT (id) DO NOTHING",
		},
		{
			name: "INSERT with ON CONFLICT DO UPDATE",
			builder: Insert("users").Columns("id", "name").
				Values([]string{"1", "'John'"}).
				OnConflict("id", false).
				DoUpdate(map[string]string{"name": "'Updated'"}),
			expected: "INSERT INTO users (id, name) VALUES (1, 'John') ON CONFLICT (id) DO UPDATE SET name = 'Updated'",
		},
		{
			name: "INSERT with ON CONFLICT ON CONSTRAINT",
			builder: Insert("users").Columns("id", "name").
				Values([]string{"1", "'John'"}).
				OnConflict("unique_id", true).
				DoNothing(),
			expected: "INSERT INTO users (id, name) VALUES (1, 'John') ON CONFLICT ON CONSTRAINT unique_id DO NOTHING",
		},
		{
			name: "INSERT with RETURNING",
			builder: Insert("users").Columns("id", "name").
				Values([]string{"1", "'John'"}).
				Returning("id", "name"),
			expected: "INSERT INTO users (id, name) VALUES (1, 'John') RETURNING id, name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.builder.ToSQL()
			if got != tt.expected {
				t.Errorf("Expected: %s, got: %s", tt.expected, got)
			}
		})
	}
}
