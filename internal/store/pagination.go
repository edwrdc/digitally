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
	Since      *string  `json:"since"`
	Until      *string  `json:"until"`
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
		rawCategories := strings.Split(categories, ",")
		cleanCategories := make([]string, 0, len(rawCategories))
		for _, cat := range rawCategories {
			cat = strings.TrimSpace(cat)
			if cat != "" {
				cleanCategories = append(cleanCategories, cat)
			}
		}
		if len(cleanCategories) > 0 {
			fq.Categories = cleanCategories
		}
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = strings.TrimSpace(search)
	}

	since := qs.Get("since")
	if since != "" {
		date, err := time.Parse(time.RFC3339, since)
		if err != nil {
			return fq, err
		}
		formattedDate := date.Format(time.RFC3339)
		fq.Since = &formattedDate
	} else {
		fq.Since = nil
	}

	until := qs.Get("until")
	if until != "" {
		date, err := time.Parse(time.RFC3339, until)
		if err != nil {
			return fq, err
		}
		formattedDate := date.Format(time.RFC3339)
		fq.Until = &formattedDate
	} else {
		fq.Until = nil
	}

	return fq, nil
}
