package backend

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/fxkr/openview/backend/cache"
	"github.com/fxkr/openview/backend/image"
	"github.com/fxkr/openview/backend/model"
	"github.com/fxkr/openview/backend/util/handler"
	"github.com/fxkr/openview/backend/util/safe"
)

type Service interface {
	Get(path safe.RelativePath) http.Handler
	GetDirectory(path safe.RelativePath, page model.Page) http.Handler
	GetImage(path safe.RelativePath) http.Handler
	GetImageThumbnail(path safe.RelativePath, size model.ThumbSize) http.Handler
}

func NewService(base safe.Path, res safe.Path, thumbnailCache cache.Cache, metadataCache cache.Cache) Service {
	return &service{base, res, thumbnailCache, metadataCache}
}

type service struct {
	base           safe.Path
	res            safe.Path
	thumbnailCache cache.Cache
	metadataCache  cache.Cache
}

// Statically assert that *service implements Service.
var _ Service = (*service)(nil)

func (s *service) Get(path safe.RelativePath) http.Handler {
	fullPath := s.base.Join(path)

	fileInfo, err := os.Stat(fullPath.String())
	if err != nil {
		return handler.StatusError(http.StatusNotFound, errors.WithStack(err))
	}

	if fileInfo.IsDir() {
		fullPath := s.res.JoinUnsafe("index.html")
		return &handler.FileHandler{Path: fullPath}
	}

	return &handler.FileHandler{Path: fullPath}
}

func (s *service) GetDirectory(path safe.RelativePath, page model.Page) http.Handler {
	fullPath := s.base.Join(path)

	var pt model.PageToken
	err := pt.UnmarshalString(page.PageToken)
	if err != nil {
		return handler.StatusError(http.StatusBadRequest, errors.WithStack(err))
	}

	fileInfos, err := ioutil.ReadDir(fullPath.String())
	if err != nil {
		return handler.StatusError(http.StatusNotFound, errors.WithStack(err))
	}

	sort.Slice(fileInfos, func(a, b int) bool {
		return model.FileInfoLessThan(fileInfos[a], fileInfos[b])
	})

	start := sort.Search(len(fileInfos), func(i int) bool {
		return pt.LessThanFileInfo(fileInfos[i])
	})

	var nextPageToken *model.PageToken
	directories := make([]model.Directory, 0)
	images := make([]model.Image, 0)
	for _, fileInfo := range fileInfos[start:] {
		relativePath := path.Join(safe.UnsafeNewRelativePath(fileInfo.Name()))

		if !pt.LessThan(relativePath.String(), fileInfo.IsDir()) {
			continue
		}

		if isImageDirectory(fileInfo) {
			directories = append(directories, model.Directory{
				Item: model.Item{
					Name:         fileInfo.Name(),
					RelativePath: relativePath,
				},
			})
			nextPageToken = &model.PageToken{
				Name:      relativePath.String(),
				Directory: true,
			}

		} else if isImage(fileInfo) {
			image, err := s.getImageData(relativePath)
			if err != nil {
				return handler.Error(err)
			}

			images = append(images, *image)
			nextPageToken = &model.PageToken{
				Name:      relativePath.String(),
				Directory: false,
			}
		}

		if len(directories)+len(images) >= page.PageSize {
			break
		}
	}

	return &handler.JSONHandler{Data: GetDirectoryResponse{
		Directory: model.Directory{
			Item: model.Item{
				Name:         path.Base(),
				RelativePath: path,
			},
		},

		Directories: directories,
		Images:      images,

		NextPageToken: nextPageToken,
	}}
}

func (s *service) GetImage(path safe.RelativePath) http.Handler {
	img, err := s.getImageData(path)
	if err != nil {
		return handler.Error(err)
	}
	return &handler.JSONHandler{Data: &GetImageResponse{*img}}
}

func (s *service) GetImageThumbnail(path safe.RelativePath, size model.ThumbSize) http.Handler {
	fullPath := s.base.Join(path)

	cacheKey := safe.NewKey("thumbnail", path.String(), size.Name)

	fileInfo, err := os.Stat(fullPath.String())
	if err != nil {
		return handler.StatusError(http.StatusNotFound, err)
	}
	if !isImage(fileInfo) {
		return handler.Status(http.StatusNotFound)
	}

	cacheVersion := s.getImageVersion(fileInfo)

	h, err := s.thumbnailCache.GetHandler(cacheKey, cacheVersion, func() (cache.Version, []byte, error) {
		bytes, err := image.RenderThumbnail(fullPath, size)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}

		return cacheVersion, bytes, nil
	}, ThumbnailContentType)

	if err != nil {
		return handler.Error(err)
	}

	return h
}

func isImageDirectory(fileInfo os.FileInfo) bool {
	if !fileInfo.Mode().IsDir() {
		return false
	}

	if strings.HasPrefix(fileInfo.Name(), ".") {
		return false
	}

	return true
}

func isImage(fileInfo os.FileInfo) bool {
	if !fileInfo.Mode().IsRegular() {
		return false
	}

	if strings.HasPrefix(fileInfo.Name(), ".") {
		return false
	}

	ext := filepath.Ext(fileInfo.Name())
	ext = strings.ToLower(ext)

	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return false
	}

	return true
}
