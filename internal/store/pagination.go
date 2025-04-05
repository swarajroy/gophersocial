package store

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" valisdate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
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

	tags := qs.Get("tags")
	if tags != "" {
		pg.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		pg.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		pg.Since = parseTime(since)
	}

	until := qs.Get("until")
	if until != "" {
		pg.Until = parseTime(until)
	}

	return pg, nil
}

func parseTime(input string) string {
	fmt.Printf("input = %s\n", input)
	conv, err := time.Parse(time.DateTime, input)
	if err != nil {
		fmt.Println("error occurred whilst parsing!")
		return ""
	}
	val := conv.Format(time.DateTime)
	fmt.Printf("val = %s", val)
	return val
}
