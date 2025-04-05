package store

import (
	"net/http"
	"strconv"
)

type PaginatedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (pg PaginatedQuery) Parse(r *http.Request) (PaginatedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pg, nil
		}
		pg.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return pg, nil
		}
		pg.Offset = o
	}

	sort := qs.Get("sort")
	if sort != "" {
		pg.Sort = sort
	}

	return pg, nil
}
