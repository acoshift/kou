package util

import (
	"net/http"
	"os"
)

// NoDirFileSystem is the file system without dir
type NoDirFileSystem struct {
	http.FileSystem
}

// Open implements FileSystem's Open
func (fs *NoDirFileSystem) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}

// NoDirFile is http.File with empty read dir
type NoDirFile struct {
	http.File
}

// Readdir implements http.File Readdir
func (f *NoDirFile) Readdir(int) ([]os.FileInfo, error) {
	return nil, os.ErrNotExist
}

// FileServer returns file server handler
func FileServer(path string) http.Handler {
	return http.FileServer(&NoDirFileSystem{http.Dir(path)})
}

// ServeFile returns handler to serve given file
func ServeFile(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	})
}
