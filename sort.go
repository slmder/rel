package rel

import "strings"

const (
	OrderAsc  = "ASC"
	OrderDesc = "DESC"
)

type ColumnOrder struct {
	Column string
	Order  string
}

type Sort []ColumnOrder

// ParseSort parses the sort query string.
// accepts format: "name,created_at:desc"
func ParseSort(query string, allowed map[string]struct{}) Sort {
	var sort Sort
	if query == "" {
		return sort
	}

	parts := strings.Split(query, ",")
	for _, part := range parts {
		if part = strings.TrimSpace(part); part == "" {
			continue
		}

		pair := strings.Split(part, ":")
		if len(pair) == 0 {
			continue
		}

		key := strings.TrimSpace(pair[0])
		if key == "" {
			continue
		}

		if len(allowed) > 0 {
			if _, exists := allowed[key]; !exists {
				continue
			}
		}

		order := OrderAsc
		if len(pair) > 1 && strings.EqualFold(strings.TrimSpace(pair[1]), "desc") {
			order = OrderDesc
		}

		sort = append(sort, ColumnOrder{key, order})
	}
	return sort
}
