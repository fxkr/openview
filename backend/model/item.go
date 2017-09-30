package model

import "github.com/fxkr/openview/backend/util/safe"

type Item struct {
	Name         string            `json:"name"`
	RelativePath safe.RelativePath `json:"relative_path"`
}
