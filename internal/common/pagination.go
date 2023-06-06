package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"math"
	"net/url"
	"strconv"
)

type PaginationData struct {
	Filters  map[string]string
	Page     int
	PageSize int
	Limit    int
	Offset   int
}

func GetPaginationParams(query url.Values) (PaginationData, error) {
	pageSize := 10
	if query.Has("page_size") {
		var err error
		pageSize, err = strconv.Atoi(query.Get("page_size"))
		query.Del("page_size")
		if err != nil {
			return PaginationData{}, errors.New("error to get 'page_size' query")
		}
		if pageSize < 1 {
			return PaginationData{}, errors.New("'page_size' can't be lower than 1")
		}
	}

	page := 1
	if query.Has("page") {
		var err error
		page, err = strconv.Atoi(query.Get("page"))
		query.Del("page")
		if err != nil {
			return PaginationData{}, errors.New("error to get 'page' query")
		}
		if page < 1 {
			return PaginationData{}, errors.New("'page' can't be lower than 1")
		}
	}

	limit := page * pageSize
	offset := (page - 1) * pageSize

	var filters = make(map[string]string, len(query))
	for k, _ := range query {
		filters[k] = query.Get(k)
	}

	pagination := PaginationData{
		Filters:  filters,
		Page:     page,
		PageSize: pageSize,
		Limit:    limit,
		Offset:   offset,
	}
	return pagination, nil
}

func PaginationJson(total int, page int, pageSize int) gin.H {
	totalPages := math.Ceil(float64(total) / float64(pageSize))

	return gin.H{
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	}
}
