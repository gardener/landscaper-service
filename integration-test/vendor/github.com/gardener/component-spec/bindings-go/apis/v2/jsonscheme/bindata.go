// Code generated by go-bindata. DO NOT EDIT.
// sources:
// ../../../../language-independent/component-descriptor-v2-schema.yaml

package jsonscheme

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  fileInfoEx
}

type fileInfoEx interface {
	os.FileInfo
	MD5Checksum() string
}

type bindataFileInfo struct {
	name        string
	size        int64
	mode        os.FileMode
	modTime     time.Time
	md5checksum string
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) MD5Checksum() string {
	return fi.md5checksum
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _LanguageIndependentComponentDescriptorV2SchemaYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x1a\x4b\x6f\xe3\xb8\xf9\xee\x5f\xf1\x61\x13\x80\xce\x24\xb2\x13\x17\x7b\x18\x5f\x82\x74\x17\x2d\x16\x2d\x36\x40\x66\xda\x43\x13\x77\x40\x4b\x9f\x6d\xa6\x12\xe9\x92\x94\x13\xcd\xe3\xbf\x17\x24\x45\x3d\x6c\x49\x7e\x25\x99\x2e\x30\x73\x98\x88\xd4\xf7\x7e\x93\xf2\x29\x8b\xc6\x40\x16\x5a\x2f\xd5\x78\x38\x9c\x53\x19\x21\x47\x39\x08\x63\x91\x46\x43\x15\x2e\x30\xa1\x6a\x18\x8a\x64\x29\x38\x72\x1d\x44\xa8\x42\xc9\x96\x5a\xc8\x60\x35\x22\xbd\x53\x07\x51\xa1\xf0\xa8\x04\x0f\xdc\xee\x40\xc8\xf9\x30\x92\x74\xa6\x87\xa3\xcb\xd1\x65\x70\x35\xca\x09\x92\x9e\x27\xc3\x04\x1f\x03\xf9\x6b\xce\x15\x7e\xf1\x7c\xe0\xd7\x82\x0f\xac\x46\x50\xa2\xcd\x18\x67\x06\x4b\x8d\x7b\x00\x09\x6a\x6a\xfe\x02\xe8\x6c\x89\x63\x20\x62\xfa\x88\xa1\x26\x76\xab\xce\xa2\xd0\x00\x4a\x0d\x2c\x7e\x44\x35\x75\x08\x12\xff\x9b\x32\x89\x91\xa3\x08\x10\x00\x71\x7c\xff\x89\x52\x31\xc1\x1d\xd4\x52\x8a\x25\x4a\xcd\x50\x79\xb8\x1a\x90\xdf\x2c\x44\x52\x5a\x32\x3e\x27\xbd\x1e\x40\x4c\xa7\x18\xb7\xca\xdb\xc0\x9e\xd3\x04\x49\xb9\x5c\xd1\x38\x45\x4b\xa9\xd0\xe6\x77\x9a\x60\x8d\xa2\x67\x67\xb6\x12\xfa\xfc\x77\xe4\x73\xbd\x18\xc3\xe8\xe7\x9f\x9d\xf4\x54\x6b\x94\xc6\x20\xff\xbe\xa7\xc1\xe7\xcb\xe0\xfd\xe0\x21\x98\x9c\xdf\x0f\x26\x66\xe9\xfe\x3b\x1f\xde\x07\xee\xdd\xf0\xd3\x60\xf2\xee\xd4\x72\x64\x11\x72\xcd\x74\x76\xa3\xb5\x64\xd3\x54\xe3\xdf\x30\x73\x8c\x13\xc6\x0b\x2e\x2d\x3c\x26\xfd\xfb\xe0\xd3\x79\xfe\xfc\xce\x6f\x9e\x5d\x3b\xd2\x12\x63\xfa\x8c\xd1\x07\x4c\x56\x28\x1d\xcd\x13\xd0\xf4\x3f\xc8\x61\x26\x45\x02\xca\xbe\x30\xc1\x04\x94\x47\x40\xa3\xc7\x54\x69\x8c\x40\x0b\xa0\x71\x2c\x9e\x80\x72\x10\xd6\xcf\x34\x86\x18\x69\xc4\xf8\x1c\xc8\x8a\x5c\x40\x42\x1f\x85\x0c\x04\x8f\xb3\x0b\x8b\x6a\xd7\x83\x84\xf1\x7c\xd7\xf3\x5a\x30\x05\x09\x52\xae\x40\x2f\x10\x66\xc2\x50\x35\x44\x9c\x31\x15\x50\x89\x86\x15\xac\x68\xcc\xa2\xba\xbc\xca\x0b\x7c\x35\x18\x0d\xfe\x54\x7d\x0e\x66\x42\x9c\x4f\xa9\xcc\xf7\x56\x55\x80\x55\x13\xc4\xd5\x60\xe4\x9f\x0a\xb0\x0a\x7c\xf1\x58\x43\xab\x1a\x7b\x35\xb9\xee\x5f\x7e\xbd\xbf\x0a\xde\x4f\x1e\xa2\x77\x67\xfd\xeb\xf1\xc3\xa0\xba\x71\x76\xdd\xbc\x15\xf4\xfb\xd7\xe3\x72\xf3\xeb\x43\x64\x7d\x74\x13\xfc\x2b\x98\xdc\x5f\x06\xef\xfd\xb3\x27\xb9\x23\xf0\x99\xe7\x78\xde\xaf\xbe\x38\xb7\x44\x6a\x3b\x16\xf2\x94\x34\xc5\x71\x53\xe8\xb5\xa6\x50\x9e\x9b\x99\xc9\x0a\x35\x86\x2f\x70\x2a\x71\x36\x06\x72\x32\xac\x14\x8e\x61\x53\x28\x13\xf8\xe6\x42\x71\x29\x14\xd3\x42\x66\xbf\x08\xae\xf1\x59\xef\x93\xad\x06\xaa\xad\x46\x58\x0a\x1d\xa5\x41\x84\xec\xae\x99\x37\x8d\xe3\xdb\x59\xc9\xa5\x51\xa3\x0d\xb1\xcb\xa2\xb1\x2e\xa7\x95\x74\x4a\x15\xfe\x43\xc6\xa4\xd8\xdb\x14\xd8\xfc\xcb\xc1\xaa\x5b\x8d\x75\xa6\x49\xc5\x66\x35\x69\x18\xa2\x52\x3b\x96\x6c\xc3\xde\x42\xc1\x4c\xc8\x1c\x15\x15\xf4\xcd\x0a\x9f\x35\x72\x53\x6f\xd5\xd9\x16\x7f\xf4\x00\xe6\x4c\x2f\xd2\xe9\x4d\x37\xef\x4e\x87\xda\xa5\xb1\x72\xc5\x6a\x76\x67\x76\x90\xc3\xfd\x36\xf2\x34\x19\xc3\x3d\x71\x02\x92\x49\xfe\x22\x67\xb4\x05\xdd\x04\x42\x37\x44\x28\x92\x84\xe9\xae\xb0\xe3\x82\xe3\x31\x76\x39\x52\xef\xdf\x05\x47\x32\x31\x82\x28\x91\xca\x10\x7f\x2d\x62\x7a\x0f\x71\x4c\x93\x2c\x16\x2b\xd7\x85\x8b\xb5\xa1\x50\x2c\x5c\x08\xb5\x08\xce\x8b\x4e\xda\x21\xf8\xee\xf5\x24\x47\xc1\x67\x2d\xe9\x6f\x39\xc0\x78\x4f\x3a\x9e\xc8\x6a\x7d\xb4\x68\x29\x02\x95\xb6\x44\x0e\x0c\xc3\x22\x06\xed\xac\xa2\x36\x50\xa9\x94\x34\x2b\x31\x99\xc6\xa4\x56\x32\x1a\x25\xb3\xb4\x3c\x52\xb5\x04\xd8\x35\xcf\x6e\x67\x55\x12\x2d\x35\xce\xe1\x91\xed\x80\xd5\x6c\xdf\x01\xdc\xcc\xad\x1e\xb8\x07\x10\xb1\x39\x2a\xfd\x61\x89\xe1\x1e\x21\xb8\xa0\x6a\x71\x13\xcf\x85\x64\x7a\x91\x94\x81\x29\x64\x42\x63\xa6\xa8\x61\xb4\xf9\xda\x4e\x73\x2d\xc1\x58\x23\xb8\xee\x04\xe7\x3e\x1f\xb6\x8d\x4c\x3a\x51\x2c\xe3\x16\x08\x93\x8a\x6c\xce\xa9\x4e\x25\xee\x69\x04\xda\xa1\xa1\x59\x25\x18\x31\xfa\xd1\xe7\xe3\xa6\xce\xf4\x68\xe1\xdd\x56\xc1\xa7\x84\xaa\xf7\x95\x8f\x0b\x74\x40\xae\xb9\x88\x99\x9d\xfa\x0a\xb5\x21\x1f\xb3\xb7\xda\xe7\xd0\x1a\xe5\x42\xac\x58\x16\xf4\xf6\x28\x4c\x35\x85\x1d\xbd\x2d\xd5\xa1\x8c\x6b\xaf\xd9\x9a\x1e\xad\x98\xb5\x78\xb0\x39\xa2\x64\x78\xe7\x9b\xcf\xd6\x2e\x4e\x4d\xa3\x42\x89\x3c\x44\x3b\xb1\x43\xbf\x3c\x4c\xc6\x22\xa4\xf1\x59\x5e\xfc\xdb\x3a\x8a\x2f\x8b\x1f\x30\xc6\x50\x0b\x79\x68\x15\x7d\x85\x8a\x56\x3d\x89\xdd\x79\x2d\x0f\xb5\x4b\x41\x69\xd7\xe3\x60\xed\x10\x58\x3d\x26\x76\x1f\x57\x1b\xce\x8e\xad\x7a\x36\xb2\xe8\xea\x94\x70\x02\x34\xd4\x29\x8d\xe3\x6c\x5c\x72\x0a\x6c\xa2\x3d\x0d\x41\x2d\x31\x64\x34\x06\x89\x06\x3e\xb4\x4c\xfe\xb8\xcd\xf5\xd5\x7a\xe4\x7a\x46\x0b\x8e\xd5\x1e\x19\x78\x4e\x3c\x8d\x2b\x43\x7c\x4b\x83\xab\x66\xbe\x3d\xe2\xb8\x74\x2b\x2b\xe4\x9e\x83\xb8\x27\xa0\x76\xbe\xb6\xc8\xe3\x11\x4e\x2c\xbe\x4d\xfa\x92\xca\x45\x7e\xfc\x4e\x95\x86\x84\xea\x70\x51\x49\x04\xb5\x31\xcf\x6d\xce\xe4\xb1\xed\x7c\x95\xad\xea\xa0\xf0\x63\xcc\x2b\xb4\x72\x45\xfb\x85\xa2\xd5\x11\x2b\x4f\x22\xce\x09\x3b\x0f\x9a\x36\x04\xc8\x05\x10\x73\x8c\x93\x9c\xc6\xdf\x7d\xec\xdc\x71\xe8\x6c\x01\x13\x21\xfb\x73\x2c\x36\x66\xce\x16\x68\xab\xfd\x5f\x58\x8c\x2a\x53\x1a\x93\x7d\x31\x6f\x9b\x98\xbd\x66\xc5\x10\x21\xfb\x2d\xa1\xf3\xa3\x0e\x8a\x76\xc9\x0c\x95\xa2\x4f\xbe\xc8\x09\xd2\xde\x9b\xcc\x99\xd2\x32\x2b\x62\xa8\xce\x66\xcb\xad\x4b\x69\xca\x1d\x15\xab\xa9\x15\x00\x89\x69\xe6\xf3\xf0\x38\x5d\x80\xe4\xe2\x10\x28\x2f\x02\x66\x6d\x43\xec\x8d\x11\xbe\x3e\x42\x98\x29\x36\xa1\x9c\xcd\x50\xe9\xf5\xf1\x75\x8d\xe9\x81\x33\xb2\xb3\x8a\x2b\xd8\x2e\x35\x9c\x04\x0a\xb4\xd8\xc2\x71\x3d\x40\x37\xd9\x39\x08\xcf\x4a\x53\x39\x47\x8d\x11\x84\x82\xeb\x62\x28\x6a\x25\xaf\xd8\xe7\x4e\x5d\xcc\x7b\x60\x1c\xa6\x99\x46\xe5\x79\x4c\x8d\xb1\xd7\xe9\xf2\x34\x99\x1a\x87\xf6\x00\x5a\x13\xf5\x88\x1c\x98\xb1\x18\xcb\xfe\x78\x6c\xc4\x34\x48\x58\x46\x8f\x67\xd5\x66\x17\xff\xbe\x6a\x0e\xd0\x0b\xaa\x81\x29\xab\xbb\x31\x3f\xe3\xf6\xdd\x4f\xe6\xa5\xfa\x09\x22\x26\xed\x10\x9e\xb5\xfa\xc3\xdb\xed\xf6\x80\xdc\x7a\x23\x83\xdd\xae\xe7\x59\x77\x70\xd6\x03\xd3\xe6\x3b\x3c\x31\xbd\xc8\x4d\x13\xa6\x52\x22\xd7\xd0\xf4\x85\xa9\xcb\x4a\xbe\xac\xde\xe5\x93\xd0\x31\x1f\x86\xaa\x13\x7f\x93\x11\x7f\xcc\x44\xdb\xfb\x88\x75\xc6\xdb\x0f\x22\x6d\x03\x45\xa5\xe5\xbe\x45\x93\x2f\xaf\xc1\x8e\xc8\xd5\xd4\xdf\x8e\x1f\xd9\xd5\x8d\x30\x85\x27\xd2\x8e\x9b\xf0\x1e\xc0\x1c\x39\x4a\x16\x7e\xc7\x5b\xec\x5c\x02\x77\x91\x9d\x2f\x7e\x24\xf5\xff\x41\x52\x97\x8e\x71\xfb\xdf\x37\xa7\x6b\x81\xfa\x16\x29\x5d\x34\xa4\x9d\x6f\xa4\xf6\xbe\x82\xda\x8c\xd1\x8d\x4f\x91\xaa\xf2\x72\x29\xc5\x8a\x45\xa5\x37\x03\x20\xb5\xbb\x84\xfa\xb5\x56\x31\xc2\xab\x1a\xfd\x1a\xc6\xb6\xb8\xdf\xfd\x56\xeb\x88\xa0\xdc\xd4\x79\xef\x18\xdb\xf8\x0a\xd2\x75\xd6\xdc\xf8\x52\x4c\xe0\xc4\x8f\x21\x71\x76\x01\x4f\x08\x82\xc7\x59\xfe\xeb\x08\x3b\xad\x0b\xee\xef\x9f\xbd\x0f\xb6\x64\xd1\xab\xe5\x4a\xee\xbe\x17\xba\x87\x58\xfb\x6c\xe8\xf1\x1b\x62\xe8\x65\x18\x6e\x12\x2e\x83\xe0\x50\xcd\x76\xf7\x7d\xf5\xee\x8e\xec\x18\x2c\xb5\x19\x73\x27\xa4\xb5\x16\x66\x6b\x49\xb3\x49\xe1\xcb\xb7\x5e\xaf\xb7\x56\x58\xaa\x55\x23\x00\x92\xa0\xfb\x7d\x55\x35\xb3\x49\xaf\x9e\xb7\xe5\xef\xb8\x1a\x05\xf2\x24\xd6\x0a\x5a\xb7\x83\x48\xf5\x53\x4d\x7d\x30\xa8\x38\xa4\xe6\x8c\xee\xcf\x1f\xa4\xf7\xbf\x00\x00\x00\xff\xff\xbc\x3c\x49\xda\x2b\x27\x00\x00")

func LanguageIndependentComponentDescriptorV2SchemaYamlBytes() ([]byte, error) {
	return bindataRead(
		_LanguageIndependentComponentDescriptorV2SchemaYaml,
		"../../../../language-independent/component-descriptor-v2-schema.yaml",
	)
}

func LanguageIndependentComponentDescriptorV2SchemaYaml() (*asset, error) {
	bytes, err := LanguageIndependentComponentDescriptorV2SchemaYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name:        "../../../../language-independent/component-descriptor-v2-schema.yaml",
		size:        10026,
		md5checksum: "",
		mode:        os.FileMode(420),
		modTime:     time.Unix(1681220532, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
// nolint: deadcode
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// AssetNames returns the names of the assets.
// nolint: deadcode
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"../../../../language-independent/component-descriptor-v2-schema.yaml": LanguageIndependentComponentDescriptorV2SchemaYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{
					Op:   "open",
					Path: name,
					Err:  os.ErrNotExist,
				}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  os.ErrNotExist,
		}
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{Func: nil, Children: map[string]*bintree{
	"..": {Func: nil, Children: map[string]*bintree{
		"..": {Func: nil, Children: map[string]*bintree{
			"..": {Func: nil, Children: map[string]*bintree{
				"..": {Func: nil, Children: map[string]*bintree{
					"language-independent": {Func: nil, Children: map[string]*bintree{
						"component-descriptor-v2-schema.yaml": {Func: LanguageIndependentComponentDescriptorV2SchemaYaml, Children: map[string]*bintree{}},
					}},
				}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}