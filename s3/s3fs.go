package s3

import (
	"errors"
	"fmt"
	"io"

	afs "github.com/qor/assetfs"
)

// AssetFileSystem AssetFS based on S3
type AssetFileSystem struct {
	paths        []string
	nameSpacedFS map[string]afs.Interface
	client       *Client
}

func NewAssetFS(config *Config) *AssetFileSystem {
	return &AssetFileSystem{
		client: New(config),
	}
}

// RegisterPath register view paths
func (fs *AssetFileSystem) RegisterPath(path string) error {
	if path == "" {
		return errors.New("path is empty")
	}

	// ignore if path already in the list
	for _, p := range fs.paths {
		if p == path {
			return nil
		}
	}

	// check if path exists.
	// TODO: implement a better approach
	list, err := fs.client.List(path)
	if err != nil {
		return fmt.Errorf("PrependPath: %+v", err)
	}

	if len(list) > 0 {
		fs.paths = append(fs.paths, path)
		return nil
	}

	return errors.New("not found")
}

// PrependPath prepend path to view paths
func (fs *AssetFileSystem) PrependPath(path string) error {
	if path == "" {
		return errors.New("path is empty")
	}

	// ignore if path already in the list
	for _, p := range fs.paths {
		if p == path {
			return nil
		}
	}

	// check if path exists.
	// TODO: implement a better approach. ObjectExists doesn't works for folders,
	// maybe get parent's list and check if folder exists there
	list, err := fs.client.List(path)
	if err != nil {
		return fmt.Errorf("PrependPath: %+v", err)
	}

	if len(list) > 0 {
		fs.paths = append([]string{path}, fs.paths...)
		return nil
	}

	return errors.New("not found")
}

// Asset get content with name from assetfs
func (fs *AssetFileSystem) Asset(name string) ([]byte, error) {
	for _, path := range fs.paths {
		key := path + "/" + name
		exists, err := fs.client.ObjectExists(key)
		if err != nil {
			return []byte{}, fmt.Errorf("checking if %s exists: %+v", key, err)
		}
		if exists {
			body, err := fs.client.Get(key)
			defer body.Close()
			if err != nil {
				return []byte{}, fmt.Errorf("getting %s: %+v", key, err)
			}
			res, err := io.ReadAll(body)
			if err != nil {
				return []byte{}, fmt.Errorf("reading %s: %+v", key, err)
			}
			return res, nil
		}
	}
	return []byte{}, fmt.Errorf("%v not found", name)
}

// Glob list matched files from assetfs
func (fs *AssetFileSystem) Glob(pattern string) (matches []string, err error) {
	// for _, pth := range fs.paths {
	// }
	return
}

// Compile compile assetfs
func (fs *AssetFileSystem) Compile() error {
	return nil
}

// NameSpace return namespaced filesystem
func (fs *AssetFileSystem) NameSpace(nameSpace string) afs.Interface {
	if fs.nameSpacedFS == nil {
		fs.nameSpacedFS = map[string]afs.Interface{}
	}
	fs.nameSpacedFS[nameSpace] = &AssetFileSystem{}
	return fs.nameSpacedFS[nameSpace]
}
