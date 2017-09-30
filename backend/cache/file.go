package cache

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dchest/safefile"
	"github.com/pkg/errors"

	"github.com/fxkr/openview/backend/util/handler"
	"github.com/fxkr/openview/backend/util/safe"
)

// FileCache is Cache implementation that stores keys as files in a directory.
type FileCache struct {
	path safe.Path
}

// Statically assert that *FileCache implements Cache.
var _ Cache = (*FileCache)(nil)

func NewFileCache(path safe.Path) (*FileCache, error) {

	stat, err := os.Stat(path.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !stat.IsDir() {
		return nil, errors.Errorf("RedisCache directory does not exist: %v", path.String())
	}

	return &FileCache{path}, nil
}

func (c *FileCache) Put(key Key, buffer []byte) error {

	// Created temporary file
	filePath := c.getFilePath(key)
	f, err := safefile.Create(filePath.String(), 0644)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	// Write to temporary file
	n, _ := f.Write(buffer)
	if n < len(buffer) {
		return errors.WithStack(io.ErrShortWrite)
	}

	// Atomically move temporary file to final location
	err = f.Commit()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *FileCache) GetBytes(key Key, filler func() ([]byte, error)) ([]byte, error) {
	file := c.getFilePath(key)

	err := c.checkFile(file)
	if err == nil {
		return ioutil.ReadFile(file.String())
	}

	value, err := filler()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = c.Put(key, value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return value, nil
}

func (c *FileCache) GetHandler(key Key, filler func() ([]byte, error), contentType string) (http.Handler, error) {
	file := c.getFilePath(key)

	err := c.checkFile(file)
	if err == nil {
		return &handler.FileHandler{Path: file}, nil // Cache hit
	}

	value, err := filler()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = c.Put(key, value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &handler.ByteHandler{
		Bytes:       value,
		ContentType: contentType,
	}, nil
}

func (c *FileCache) getFileName(key Key) safe.RelativePath {
	unsafeFilename := key.String()
	safeFilename := base64.URLEncoding.EncodeToString([]byte(unsafeFilename))
	return safe.UnsafeNewRelativePath(safeFilename)
}

func (c *FileCache) getFilePath(key Key) safe.Path {
	return c.path.Join(c.getFileName(key))
}

func (c *FileCache) checkFile(file safe.Path) error {
	stat, err := os.Stat(file.String())
	if err != nil {
		return errors.WithStack(err)
	}
	if !stat.Mode().IsRegular() {
		return errors.Errorf("Corrupt cache: not a regular file: %s", file.String())
	}
	return nil
}

func (c *FileCache) Close() {

}
