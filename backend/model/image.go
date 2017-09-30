package model

type Image struct {
	Item

	Width  uint `json:"width"`
	Height uint `json:"height"`
}
