package image

import (
	"github.com/pkg/errors"
	"gopkg.in/gographics/imagick.v2/imagick"

	"github.com/fxkr/openview/backend/model"
	"github.com/fxkr/openview/backend/util/safe"
)

func GetImageData(fullPath safe.Path) (*model.Image, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImage(fullPath.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	orientation := mw.GetImageOrientation()
	if orientation == imagick.ORIENTATION_LEFT_TOP ||
		orientation == imagick.ORIENTATION_RIGHT_TOP ||
		orientation == imagick.ORIENTATION_RIGHT_BOTTOM ||
		orientation == imagick.ORIENTATION_LEFT_BOTTOM {
		width, height = height, width
	}

	return &model.Image{
		Item: model.Item{},

		Width:  width,
		Height: height,
	}, nil
}
