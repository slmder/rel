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
		pair := strings.Split(part, ":")
		key := strings.TrimSpace(pair[0])
		if key == "" {
			continue
		}
		order := OrderAsc
		switch {
		case len(pair) > 1 && strings.TrimSpace(strings.ToLower(pair[1])) == "desc":
			order = OrderDesc
		default:
			order = OrderAsc
		}
		if len(allowed) > 0 {
			if _, exists := allowed[key]; !exists {
				continue
			}
		}
		sort = append(sort, ColumnOrder{key, order})
	}
	return sort
}
