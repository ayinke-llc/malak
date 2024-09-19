package malak

import (
	"net/http"
	"strconv"

	"github.com/ayinke-llc/malak/internal/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

const (
	defaultNumOfItemPerPage = 8
)

type PaginatedResultMetadata struct {
	Total int64
}

type Paginator struct {
	PerPage int64
	Page    int64
}

func (p Paginator) OTELAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Int64("per_page", p.PerPage),
		attribute.Int64("page", p.Page),
	}
}

func (p Paginator) Offset() int64 {
	if p.Page <= 0 {
		return 0
	}

	return (p.Page - 1) * p.PerPage
}

func PaginatorFromRequest(r *http.Request) Paginator {
	defaultPage := 1

	p := Paginator{
		Page:    int64(defaultPage),
		PerPage: defaultNumOfItemPerPage,
	}

	if !util.IsStringEmpty(r.URL.Query().Get("page")) {
		currentPage := r.URL.Query().Get("page")

		var err error

		dd, err := strconv.Atoi(currentPage)
		if err != nil || p.Page <= 0 {
			return p
		}

		p.Page = int64(dd)
	}

	if !util.IsStringEmpty(r.URL.Query().Get("per_page")) {
		perPage := r.URL.Query().Get("per_page")
		var err error

		dd, err := strconv.Atoi(perPage)
		if err != nil || p.PerPage <= 0 {
			return p
		}

		p.PerPage = int64(dd)
	}

	return p
}
