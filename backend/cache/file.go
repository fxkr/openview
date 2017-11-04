package cache

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dchest/safefile"
	"github.com/pkg/errors"
	"github.com/pkg/xattr"

	"github.com/fxkr/openview/backend/util/handler"
	"github.com/fxkr/openview/backend/util/safe"
)

// FileCache is Cache implementation that stores keys as files in a directory.
//
// Metadata (currently only the version, used for expiration) is stored in extended attributes,
// so the filesystem needs to support these. Nearly all Linux filesystems do.
type FileCache struct {
	path safe.Path
}

// Statically assert that *FileCache implements Cache.
var _ Cache = (*FileCache)(nil)

// versionXattr is the name of the extended attribute used to store the cache item version.
//
// Names of extended attributes used by userspace tools must start with "user.".
const versionXattr = "user.openview.cache-version"

func NewFileCache(path safe.Path) (*FileCache, error) {

	stat, err := os.Stat(path.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !stat.IsDir() {
		return nil, errors.Errorf("Cache directory does not exist: %v", path.String())
	}
	if !xattr.Supported(path.String()) {
		return nil, errors.Errorf("Cache directorie's file system does not support extended attributes: %v", path.String())
	}

	return &FileCache{path}, nil
}

func (c *FileCache) Put(key Key, version Version, buffer []byte) error {

	// Created temporary file
	// (Same directory, so move will be atomic and xattrs won't get lost.)
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

	// Store version as extended attribute
	err = xattr.Set(f.Name(), versionXattr, []byte(version.String()))
	if err != nil {
		f.Close() // Deletes the temporary file
		return errors.WithStack(err)
	}

	// Atomically move temporary file to final location
	err = f.Commit()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *FileCache) GetBytes(key Key, version Version, filler func() (Version, []byte, error)) ([]byte, error) {
	file := c.getFilePath(key)

	err := c.checkFile(file, version)
	if err == nil {
		// Cache hit
		return ioutil.ReadFile(file.String())
	}

	version, value, err := filler()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = c.Put(key, version, value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return value, nil
}

func (c *FileCache) GetHandler(key Key, version Version, filler func() (Version, []byte, error), contentType string) (http.Handler, error) {
	file := c.getFilePath(key)

	err := c.checkFile(file, version)
	if err == nil {
		// Cache hit
		return &handler.FileHandler{Path: file}, nil // Cache hit
	}

	version, value, err := filler()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = c.Put(key, version, value)
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

func (c *FileCache) checkFile(file safe.Path, requestedVersion Version) error {
	stat, err := os.Stat(file.String())
	if err != nil {
		return errors.WithStack(err)
	}
	if !stat.Mode().IsRegular() {
		return errors.Errorf("Corrupt cache: not a regular file: %s", file.String())
	}
	cachedVersion, err := xattr.Get(file.String(), versionXattr)
	if err != nil {
		return errors.Errorf("Failed to read extended attributes: %s", file.String())
	}
	if !bytes.Equal(cachedVersion, []byte(requestedVersion.String())) {
		return errors.Errorf("Outdated cache item: %s", file.String())
	}
	return nil
}

func (c *FileCache) Close() {

}
