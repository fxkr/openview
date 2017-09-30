package model

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// PageToken represents a position in a list of items.
//
// The frontend should treat is as an opaque value.
//
// The item order is directories first, then alphabetic by name.
// A nil page token signifies that there are no more items available.
type PageToken struct {
	Name      string
	Directory bool
}

// pageToken is used to marshal PageToken to JSON.
type pageToken struct {
	Name      string `json:"n"`
	Directory bool   `json:"d"`
}

func (pt *PageToken) MarshalJSON() ([]byte, error) {
	if pt == nil {
		return []byte("null"), nil
	}

	s, err := pt.MarshalString()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return json.Marshal(s)
}

func (pt *PageToken) UnmarshalJSON(data []byte) error {
	var base64String string

	err := json.Unmarshal(data, &base64String)
	if err != nil {
		return errors.WithStack(err)
	}

	return pt.UnmarshalString(base64String)
}

func (pt *PageToken) UnmarshalString(data string) error {
	if data == "" {
		pt.Name = ""
		pt.Directory = true
		return nil
	}

	jsonString, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return errors.WithStack(err)
	}

	var rawPageToken pageToken
	err = json.Unmarshal([]byte(jsonString), &rawPageToken)
	if err != nil {
		return errors.WithStack(err)
	}

	pt.Name = rawPageToken.Name
	pt.Directory = rawPageToken.Directory
	return nil
}

func (pt *PageToken) MarshalString() (string, error) {
	if pt == nil {
		return "", nil
	}

	b, err := json.Marshal(&pageToken{
		Name:      pt.Name,
		Directory: pt.Directory,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (pt *PageToken) LessThan(name string, isDir bool) bool {
	if pt.Directory != isDir {
		return pt.Directory
	}
	return pt.Name < name
}

func FileInfoLessThan(a, b os.FileInfo) bool {
	if a.IsDir() != b.IsDir() { // directories before files
		return a.IsDir()
	} else { // otherwise, alphabetically
		return a.Name() < b.Name()
	}
}

func (pt *PageToken) LessThanFileInfo(fi os.FileInfo) bool {
	return pt.LessThan(fi.Name(), fi.IsDir())
}
