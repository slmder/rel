package rel

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestConditions(t *testing.T) {
	tests := []struct {
		name     string
		cond     Cond
		expected []string
		args     []any
	}{
		{
			name:     "Eq",
			cond:     Cond{Eq("name", "Alice")},
			expected: []string{"\"name\" = $1"},
			args:     []any{"Alice"},
		},
		{
			name:     "Any",
			cond:     Cond{Any("name", "Alice")},
			expected: []string{"\"name\" = ANY $1"},
			args:     []any{pq.Array([]any{"Alice"})},
		},
		{
			name:     "Any",
			cond:     Cond{Any("name", Identifier("user_names"))},
			expected: []string{"\"name\" = ANY \"user_names\""},
		},
		{
			name:     "NotAll",
			cond:     Cond{NotAll("name", "Alice")},
			expected: []string{"\"name\" <> ALL $1"},
			args:     []any{pq.Array([]any{"Alice"})},
		},
		{
			name:     "NotAll",
			cond:     Cond{NotAll("name", Identifier("user_names"))},
			expected: []string{"\"name\" <> ALL \"user_names\""},
		},
		{
			name:     "Contains",
			cond:     Cond{Contains("name", 1, 2, 3)},
			expected: []string{"\"name\"  @> $1"},
			args:     []any{pq.Array([]any{1, 2, 3})},
		},
		{
			name:     "Eq",
			cond:     Cond{Eq("name", Identifier("user_name"))},
			expected: []string{"\"name\" = \"user_name\""},
		},
		{
			name:     "Neq",
			cond:     Cond{Neq("age", 30)},
			expected: []string{"\"age\" <> $1"},
			args:     []any{30},
		},
		{
			name:     "Neq",
			cond:     Cond{Neq("age", Identifier("user_age"))},
			expected: []string{"\"age\" <> \"user_age\""},
		},
		{
			name:     "In",
			cond:     Cond{In("id", 1, 2, 3)},
			expected: []string{"\"id\" IN ($1,$2,$3)"},
			args:     []any{1, 2, 3},
		},
		{
			name:     "NotIn",
			cond:     Cond{NotIn("status", "active", "inactive")},
			expected: []string{"\"status\" NOT IN ($1,$2)"},
			args:     []any{"active", "inactive"},
		},
		{
			name:     "IsNull",
			cond:     Cond{IsNull("deleted_at")},
			expected: []string{"\"deleted_at\" IS NULL"},
			args:     nil,
		},
		{
			name:     "NotNull",
			cond:     Cond{NotNull("deleted_at")},
			expected: []string{"\"deleted_at\" IS NOT NULL"},
			args:     nil,
		},
		{
			name:     "Gt",
			cond:     Cond{Gt("score", 50)},
			expected: []string{"\"score\" > $1"},
			args:     []any{50},
		},
		{
			name:     "Gt",
			cond:     Cond{Gt("score", Identifier("min_score"))},
			expected: []string{"\"score\" > \"min_score\""},
		},
		{
			name:     "Gte",
			cond:     Cond{Gte("score", 75)},
			expected: []string{"\"score\" >= $1"},
			args:     []any{75},
		},
		{
			name:     "Gte",
			cond:     Cond{Gte("score", Identifier("max_score"))},
			expected: []string{"\"score\" >= \"max_score\""},
		},
		{
			name:     "Lt",
			cond:     Cond{Lt("score", 40)},
			expected: []string{"\"score\" < $1"},
			args:     []any{40},
		},
		{
			name:     "Lt",
			cond:     Cond{Lt("score", Identifier("max_score"))},
			expected: []string{"\"score\" < \"max_score\""},
		},
		{
			name:     "Lte",
			cond:     Cond{Lte("score", 100)},
			expected: []string{"\"score\" <= $1"},
			args:     []any{100},
		},
		{
			name:     "Lte",
			cond:     Cond{Lte("score", Identifier("max_score"))},
			expected: []string{"\"score\" <= \"max_score\""},
		},
		{
			name:     "Between",
			cond:     Cond{Between("created_at", "2023-01-01", "2023-12-31")},
			expected: []string{"\"created_at\" BETWEEN $1 AND $2"},
			args:     []any{"2023-01-01", "2023-12-31"},
		},
		{
			name:     "Between",
			cond:     Cond{Between("created_at", Identifier("from"), Identifier("to"))},
			expected: []string{"\"created_at\" BETWEEN \"from\" AND \"to\""},
			args:     nil,
		},
		{
			name:     "Like",
			cond:     Cond{Like("username", "john")},
			expected: []string{"\"username\" LIKE $1"},
			args:     []any{"%john%"},
		},
		{
			name:     "LikeLower",
			cond:     Cond{LikeLower("username", "John")},
			expected: []string{"LOWER(\"username\") LIKE $1"},
			args:     []any{"%john%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, expr := tt.cond.Split()
			assert.Equal(t, tt.expected, expr)
			assert.Equal(t, tt.args, args)
		})
	}
}
