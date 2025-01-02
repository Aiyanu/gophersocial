package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int        `json:"limit" validate:"gte=1,lte=20"`
	Offset int        `json:"offset" validate:"gte=0"`
	Sort   string     `json:"sort" validate:"oneof=asc desc"`
	Tags   []string   `json:"tags" validate:"max=5,dive,required"`
	Search string     `json:"search" validate:"max=100"`
	Since  *time.Time `json:"since"`
	Until  *time.Time `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	}
	if offset := qs.Get("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = o
	}
	if sort := qs.Get("sort"); sort != "" {
		fq.Sort = sort
	}
	if tags := qs.Get("tags"); tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}
	if search := qs.Get("search"); search != "" {
		fq.Search = search
	}
	if since := qs.Get("since"); since != "" {
		t, err := time.Parse(time.RFC3339, since)
		if err != nil {
			return fq, err
		}
		fq.Since = &t
	}
	if until := qs.Get("until"); until != "" {
		t, err := time.Parse(time.RFC3339, until)
		if err != nil {
			return fq, err
		}
		fq.Until = &t
	}

	return fq, nil
}
