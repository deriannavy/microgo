package store

import (
	"net/http"
	"strconv"
)

type PaginatedTransactionQuery struct {
	Limit  int    `json:"l" validate:"gte=1,max=50"`
	Offset int    `json:"o" validate:"gte=0"`
	Sort   string `json:"s" validate:"omitempty,oneof=asc desc"`
}

func (fq PaginatedTransactionQuery) Parse(r *http.Request) (PaginatedTransactionQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("l")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	}
	offset := qs.Get("o")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = o
	}
	sort := qs.Get("s")
	if sort != "" {
		fq.Sort = sort
	}
	return fq, nil
}
