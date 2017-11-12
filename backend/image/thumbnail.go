package image

import (
	"github.com/pkg/errors"
	"gopkg.in/gographics/imagick.v2/imagick"

	"github.com/fxkr/openview/backend/model"
	"github.com/fxkr/openview/backend/util/safe"
)

func RenderThumbnail(fullPath safe.Path, size model.ThumbSize) ([]byte, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err := mw.ReadImage(fullPath.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	maxOldSize := width
	if height > maxOldSize {
		maxOldSize = height
	}

	if maxOldSize > size.Pixel {
		factor := float32(size.Pixel) / float32(maxOldSize)
		width = uint(float32(width) * factor)
		height = uint(float32(height) * factor)
	}

	err = mw.ResizeImage(width, height, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = mw.AutoOrientImage()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = mw.SetImageCompressionQuality(95)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = mw.SetImageFormat("JPG")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	mw.ResetIterator()

	return mw.GetImageBlob(), nil
}
