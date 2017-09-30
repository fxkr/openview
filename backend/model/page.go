package model

import (
	"strconv"

	"github.com/pkg/errors"
)

const (
	MAX_PAGE_SIZE     = 100
	DEFAULT_PAGE_SIZE = 25
)

type Page struct {
	PageToken string `json:"page_token"`
	PageSize  int    `json:"page_size"`
}

func NewPage(pageToken string, pageSize string) (Page, error) {
	if pageSize == "" {
		return Page{pageToken, DEFAULT_PAGE_SIZE}, nil
	}

	pageSizeUint, err := strconv.ParseUint(pageSize, 10, 16)
	if err != nil {
		return Page{}, errors.WithStack(err)
	}

	if pageSizeUint <= 0 {
		pageSizeUint = DEFAULT_PAGE_SIZE
	} else if pageSizeUint > MAX_PAGE_SIZE {
		pageSizeUint = MAX_PAGE_SIZE
	}

	return Page{pageToken, int(pageSizeUint)}, nil
}
