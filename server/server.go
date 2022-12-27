package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/MAA-Contest-Tester/search/database"
)

var client *database.SearchClient
var logger = log.New(os.Stderr, "[HTTP Server]  ", 0)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "API Reached!")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	if q != "" {
		res, err := client.Search(q)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprint(err))
		} else {
			io.WriteString(w, res)
		}
	} else {
		io.WriteString(w, "[]")
	}
}

func InitServer(path *string) *http.ServeMux {
	mux := http.NewServeMux()
	client = database.Client()
	if path != nil {
		fileserver := http.FileServer(http.Dir(*path))
		mux.Handle("/", fileserver)
	} else {
		mux.HandleFunc("/", indexHandler)
	}
	mux.HandleFunc("/search", searchHandler)
	return mux
}
