package model

import (
	"github.com/pkg/errors"
)

type ThumbSize struct {
	Name  string `json:"name"`
	Pixel uint   `json:"pixel"`
}

var (
	ThumbSizes = map[string]ThumbSize{
		"":     ThumbSize{"800", 800}, // "default"
		"100":  ThumbSize{"100", 100},
		"240":  ThumbSize{"240", 240},
		"360":  ThumbSize{"360", 360},
		"500":  ThumbSize{"500", 500},
		"800":  ThumbSize{"800", 800},
		"1024": ThumbSize{"1024", 1024},
		"1600": ThumbSize{"1600", 1600},
		"2048": ThumbSize{"2048", 2048},
	}
)

func NewThumbSize(s string) (ThumbSize, error) {
	result, ok := ThumbSizes[s]
	if !ok {
		return ThumbSize{}, errors.Errorf("Bad thumbnail size: %v", s)
	}
	return result, nil
}
