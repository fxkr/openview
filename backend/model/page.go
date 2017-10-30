package model

import (
	"strconv"

	"github.com/pkg/errors"
)

const (
	MaxPageSize     = 100
	DefaultPageSize = 25
)

type Page struct {
	PageToken string `json:"page_token"`
	PageSize  int    `json:"page_size"`
}

func NewPage(pageToken string, pageSize string) (Page, error) {
	if pageSize == "" {
		return Page{pageToken, DefaultPageSize}, nil
	}

	pageSizeUint, err := strconv.ParseUint(pageSize, 10, 16)
	if err != nil {
		return Page{}, errors.WithStack(err)
	}

	if pageSizeUint <= 0 {
		pageSizeUint = DefaultPageSize
	} else if pageSizeUint > MaxPageSize {
		pageSizeUint = MaxPageSize
	}

	return Page{pageToken, int(pageSizeUint)}, nil
}
