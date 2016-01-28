package backends

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const DEFAULT_FILESYSTEM_PATH = "/var/db/hoarder"

// implementation of Driver
type Filesystem struct {
	Path string // path to the local database (default)
}

// Init ensures the database exists before trying to do any operations on it
func (d *Filesystem) Init() error {
	//
	if d.Path == "" {
		d.Path = DEFAULT_FILESYSTEM_PATH
	}

	//
	return os.MkdirAll(d.Path, 0755)
}

// List returns a list of files, and some info, currently stored
func (d Filesystem) List() ([]FileInfo, error) {
	//
	files, err := ioutil.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}

	//
	info := []FileInfo{}
	for _, fi := range files {
		info = append(info, FileInfo{Name: fi.Name(), Size: fi.Size(), ModTime: fi.ModTime().UTC()})
	}

	//
	return info, nil
}

// Read reads a file and returns the contents
func (d Filesystem) Read(key string) (io.Reader, error) {
	//
	f, err := os.Open(filepath.Join(d.Path, key))
	if err != nil {
		return nil, err
	}

	//
	return f, nil
}

// Remove removes a file
func (d Filesystem) Remove(key string) error {
	//
	return os.RemoveAll(filepath.Join(d.Path, key))
}

// Stat returns information about a file
func (d Filesystem) Stat(key string) (FileInfo, error) {
	//
	fi, err := os.Stat(filepath.Join(d.Path, key))
	if err != nil {
		return FileInfo{}, err
	}

	//
	return FileInfo{Name: fi.Name(), Size: fi.Size(), ModTime: fi.ModTime().UTC()}, nil
}

// Write writes data do a file
func (d Filesystem) Write(key string, r io.Reader) error {
	// read the entire contents of the reader
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	// create/truncate a file and write the contents to it
	return ioutil.WriteFile(filepath.Join(d.Path, key), b, 0644)
}