package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginationFeedQuery struct {
	Limit      int      `json:"limit" validate:"gte=1,lte=20"`
	Offset     int      `json:"offset" validate:"gte=0"`
	Sort       string   `json:"sort" validate:"oneof=asc desc"`
	Categories []string `json:"categories" validate:"max=5"`
	Search     string   `json:"search" validate:"max=100"`
	Since      string   `json:"since"`
	Until      string   `json:"until"`
}

func (fq PaginationFeedQuery) Parse(r *http.Request) (PaginationFeedQuery, error) {

	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = o
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	categories := qs.Get("categories")
	if categories != "" {
		fq.Categories = strings.Split(categories, ",")
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		date, err := parseDate(since)
		if err != nil {
			return fq, nil
		}
		fq.Since = date.Format(time.RFC3339)
	}

	until := qs.Get("until")
	if until != "" {
		date, err := parseDate(until)
		if err != nil {
			return fq, nil
		}
		fq.Until = date.Format(time.RFC3339)
	}

	return fq, nil
}

func parseDate(date string) (time.Time, error) {
	return time.Parse(time.DateTime, date)
}
