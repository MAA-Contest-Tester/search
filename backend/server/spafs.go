package server

import (
	"log"
	"net/http"
	"os"
)

type SpaFS struct {
	root http.FileSystem;
	fallback string
}

func (f *SpaFS) Open(name string) (http.File, error) {
	file, err := f.root.Open(name)
	if os.IsNotExist(err) {
		return f.root.Open(f.fallback)
	}
	log.Println(name, file)
	return file, err
}
