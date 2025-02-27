package qbuilder

import "testing"

func TestSelectBuilder_ToSQL(t *testing.T) {
	tests := []struct {
		name     string
		builder  *SelectBuilder
		expected string
	}{
		{
			name:     "Simple SELECT *",
			builder:  Select(),
			expected: "SELECT * FROM ",
		},
		{
			name:     "SELECT columns",
			builder:  Select("id", "name").From("users"),
			expected: "SELECT id, name FROM users",
		},
		{
			name:     "SELECT DISTINCT",
			builder:  Select("name").From("users").Distinct(true),
			expected: "SELECT DISTINCT name FROM users",
		},
		{
			name:     "SELECT with WHERE",
			builder:  Select("id").From("users").Where("age > %d", 18),
			expected: "SELECT id FROM users WHERE age > 18",
		},
		{
			name:     "SELECT with multiple WHERE",
			builder:  Select("id").From("users").Where("age > %d", 18).AndWhere("status = %s", "'active'"),
			expected: "SELECT id FROM users WHERE age > 18 AND status = 'active'",
		},
		{
			name:     "SELECT with JOIN",
			builder:  Select("u.id", "o.total").From("users", "u").LeftJoin("orders", "o", "u.id = o.user_id"),
			expected: "SELECT u.id, o.total FROM users AS u LEFT JOIN orders AS o ON u.id = o.user_id",
		},
		{
			name:     "SELECT with GROUP BY",
			builder:  Select("category", "COUNT(*)").From("products").GroupBy("category"),
			expected: "SELECT category, COUNT(*) FROM products GROUP BY category",
		},
		{
			name:     "SELECT with HAVING",
			builder:  Select("category", "COUNT(*)").From("products").GroupBy("category").Having("COUNT(*) > 10"),
			expected: "SELECT category, COUNT(*) FROM products GROUP BY category HAVING COUNT(*) > 10",
		},
		{
			name:     "SELECT with ORDER BY",
			builder:  Select("id").From("users").OrderBy("created_at", OrderASC),
			expected: "SELECT id FROM users ORDER BY created_at ASC",
		},
		{
			name:     "SELECT with LIMIT and OFFSET",
			builder:  Select("id").From("users").Limit(10).Offset(5),
			expected: "SELECT id FROM users LIMIT 10 OFFSET 5",
		},
		{
			name: "SELECT with UNION",
			builder: Select("id").From("users").
				Union(Select("id").From("admins").ToSQL()),
			expected: "SELECT id FROM users UNION (SELECT id FROM admins)",
		},
		{
			name: "SELECT with UNION ALL",
			builder: Select("id").From("users").
				UnionAll(Select("id").From("admins").ToSQL()),
			expected: "SELECT id FROM users UNION ALL (SELECT id FROM admins)",
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
