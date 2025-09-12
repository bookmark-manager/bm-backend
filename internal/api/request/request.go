package request

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

const (
	DefaultPerpage = math.MaxInt32
	DefaultPage    = 1
)

type Request struct {
	Title string `json:"title" validate:"required"`
	URL   string `json:"url" validate:"required,url"`
}

type ListOptions struct {
	Perpage int
	Page    int
	Search  string
}

func (p *ListOptions) Offset() int {
	return (p.Page - 1) * p.Perpage
}

func ParseListOptions(r *http.Request) (*ListOptions, error) {
	perPage := r.URL.Query().Get("per_page")
	page := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")

	opts := &ListOptions{
		Perpage: DefaultPerpage,
		Page:    DefaultPage,
		Search:  search,
	}

	parsedLimit, err := strconv.Atoi(perPage)
	if err != nil {
		return opts, fmt.Errorf("failed to convert limit to integer: %w", err)
	}
	opts.Perpage = parsedLimit

	parsedPage, err := strconv.Atoi(page)
	if err != nil {
		return opts, fmt.Errorf("failed to convert page to integer:%w", err)
	}
	opts.Page = parsedPage

	return opts, nil
}
