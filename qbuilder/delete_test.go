package qbuilder

import "testing"

func TestDeleteBuilder_ToSQL(t *testing.T) {
	tests := []struct {
		name     string
		builder  *DeleteBuilder
		expected string
	}{
		{
			name:     "Basic delete",
			builder:  Delete("users"),
			expected: "DELETE FROM users",
		},
		{
			name:     "Delete with alias",
			builder:  Delete("users").As("u"),
			expected: "DELETE FROM users AS u",
		},
		{
			name:     "Delete with WHERE",
			builder:  Delete("users").Where("id = 1"),
			expected: "DELETE FROM users WHERE id = 1",
		},
		{
			name:     "Delete with multiple WHERE conditions",
			builder:  Delete("users").Where("id = 1").AndWhere("name = 'John'"),
			expected: "DELETE FROM users WHERE id = 1 AND name = 'John'",
		},
		{
			name:     "Delete with USING",
			builder:  Delete("users").Using("orders"),
			expected: "DELETE FROM users USING orders",
		},
		{
			name:     "Delete with multiple USING",
			builder:  Delete("users").Using("orders").AndUsing("payments"),
			expected: "DELETE FROM users USING orders, payments",
		},
		{
			name:     "Delete with RETURNING",
			builder:  Delete("users").Returning("id", "name"),
			expected: "DELETE FROM users RETURNING id, name",
		},
		{
			name:     "Complex delete with alias, USING, WHERE, RETURNING",
			builder:  Delete("users").As("u").Using("orders").Where("u.id = orders.user_id").Returning("u.id"),
			expected: "DELETE FROM users AS u USING orders WHERE u.id = orders.user_id RETURNING u.id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.builder.ToSQL()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
