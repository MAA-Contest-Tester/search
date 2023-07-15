package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MAA-Contest-Tester/search/backend/database"
)

var client database.SearchClient

func indexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "API Reached!")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	w.Header().Add("Content-Type", "application/json")
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
	client = database.InitMeiliSearchClient()
	if path != nil {
		fileserver := http.FileServer(http.Dir(*path))
		mux.Handle("/", fileserver)
	} else {
		mux.HandleFunc("/", indexHandler)
	}
	mux.HandleFunc("/search", searchHandler)
	return mux
}
