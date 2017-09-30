package backend

import (
	"github.com/fxkr/openview/backend/model"
)

type GetImageResponse struct {
	model.Image
}

type GetDirectoryResponse struct {
	model.Directory

	Directories []model.Directory `json:"directories"`
	Images      []model.Image     `json:"images"`

	NextPageToken *model.PageToken `json:"next_page_token"`
}
