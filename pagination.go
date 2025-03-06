package rel

import (
	"strconv"
)

var (
	PaginationDefaultMaxLimit uint32 = 200
)

// Pagination represents the pagination options
type Pagination struct {
	Limit uint32
	Page  uint32
}

// ComputeOffset returns the offset
func (p Pagination) ComputeOffset() uint32 {
	if p.Page <= 1 {
		return 0
	}

	return (p.Page - 1) * p.Limit
}

// ComputeLimit returns the limit
func (p Pagination) ComputeLimit() uint32 {
	if p.Limit < 1 {
		return 0
	}

	return p.Limit
}

// ParsePagination parses the Page and Limit from the request query string
func ParsePagination(page, limit string) Pagination {
	pag, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		pag = 1
	}
	lim, err := strconv.ParseUint(limit, 10, 32)
	if err != nil {
		lim = 0
	}
	return Pagination{
		Limit: uint32(lim),
		Page:  uint32(pag),
	}
}
