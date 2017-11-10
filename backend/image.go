package backend

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/gographics/imagick.v2/imagick"

	"github.com/fxkr/openview/backend/cache"
	"github.com/fxkr/openview/backend/model"
	"github.com/fxkr/openview/backend/util/handler"
	"github.com/fxkr/openview/backend/util/safe"
)

const (
	ThumbnailContentType = "image/jpeg"
)

func Initialize() {
	imagick.Initialize()
}

func Terminate() {
	imagick.Terminate()
}

type ThumbnailOptions struct {
	Width  *int
	Height *int
}

func (s *service) getImageVersion(fileInfo os.FileInfo) cache.Version {
	return safe.NewKey(fileInfo.ModTime(), fileInfo.Size())
}

func (s *service) getImageData(path safe.RelativePath) (*model.Image, error) {
	cacheKey := safe.NewKey("imagemeta", path.String())

	fullPath := s.base.Join(path)

	fileInfo, err := os.Stat(fullPath.String())
	if err != nil {
		return nil, handler.StatusError(http.StatusNotFound, errors.WithStack(err))
	}

	cacheVersion := s.getImageVersion(fileInfo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resultBuf, err := s.metadataCache.GetBytes(cacheKey, cacheVersion, func() (cache.Version, []byte, error) {
		if !isImage(fileInfo) {
			return nil, nil, handler.StatusError(http.StatusNotFound, errors.WithStack(err))
		}

		mw := imagick.NewMagickWand()
		defer mw.Destroy()
		err = mw.ReadImage(fullPath.String())
		if err != nil {
			return nil, nil, handler.StatusError(http.StatusInternalServerError, errors.WithStack(err))
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

		value := model.Image{
			Item:   model.Item{Name: path.Base(), RelativePath: path},
			Width:  width,
			Height: height,
		}

		buf, err := json.Marshal(value)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}

		return cacheVersion, buf, nil
	})

	if err != nil {
		return &model.Image{}, err
	}

	var result model.Image
	err = json.Unmarshal(resultBuf, &result)
	if err != nil {
		return &model.Image{}, err
	}

	return &result, nil
}

func (s *service) getThumbnail(path safe.RelativePath, size model.ThumbSize) ([]byte, error) {
	buf, err := s.renderThumbnail(path, size)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buf, nil
}

func (s *service) renderThumbnail(path safe.RelativePath, size model.ThumbSize) ([]byte, error) {
	fullPath := s.base.Join(path)

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
