package rel

import (
	"strconv"
)

var (
	PaginationDefaultMaxLimit uint32 = 200
)

// Pagination represents the pagination options
type Pagination struct {
	lim uint32
	pag uint32
}

// Limit sets the limit
func (p Pagination) Limit() uint32 {
	if p.lim <= 0 {
		return 0
	}

	if p.lim > PaginationDefaultMaxLimit {
		return PaginationDefaultMaxLimit
	}

	return p.lim
}

// Offset returns the offset
func (p Pagination) Offset() uint32 {
	if p.pag <= 1 {
		return 0
	}

	return (p.pag - 1) * p.Limit()
}

// ParsePagination parses the pag and lim from the request query string
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
		lim: uint32(lim),
		pag: uint32(pag),
	}
}
