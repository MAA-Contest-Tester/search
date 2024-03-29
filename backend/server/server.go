package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/MAA-Contest-Tester/search/backend/database"
	"github.com/MAA-Contest-Tester/search/backend/scrape"
)

var client database.SearchClient

func indexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "API Reached!")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	o := r.URL.Query().Get("offset")
	offset, _ := strconv.Atoi(o)
	w.Header().Add("Content-Type", "application/json")
	if q != "" {
		res, err := client.SearchProblems(q, offset)
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

type handoutRequest struct {
	ProblemIds []string `json:"ids"`
}

func handoutHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Header.Get("Content-Type")) != "application/json" {
		msg := "Content-Type header is not application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	request := handoutRequest{}
	err := dec.Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
	}
	problems := make([]*scrape.Problem, len(request.ProblemIds))
	for index, id := range request.ProblemIds {
		res, err := client.GetById(id)
		if err != nil {
			problems[index] = nil
		} else {
			err := json.Unmarshal([]byte(res), &problems[index])
			if err != nil {
				panic(err)
			}
		}
	}
	encoded, err := json.Marshal(problems)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(encoded)
}

var routes []string = []string{"/", "/meta", "/handout"}

func InitServer(path *string, meta *scrape.Meta) *http.ServeMux {
	mux := http.NewServeMux()
	client = database.InitMeiliSearchClient()
	if path != nil {
		fs := SpaFS {
			root: http.Dir(*path),
			fallback: "/index.html",
		}
		fileserver := http.FileServer(&fs)
		mux.Handle("/", fileserver)
		for _, route := range routes {
			http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, "index.html")
			})
		}
	} else {
		mux.HandleFunc("/", indexHandler)
	}
	mux.HandleFunc("/backend/search", searchHandler)
	mux.HandleFunc("/backend/handout", handoutHandler)
	mux.HandleFunc("/backend/meta", func(w http.ResponseWriter, r *http.Request) {
		encoded, err := json.Marshal(meta)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(encoded)
	})
	return mux
}
