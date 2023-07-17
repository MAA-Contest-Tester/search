package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/MAA-Contest-Tester/search/backend/database"
	"github.com/MAA-Contest-Tester/search/backend/scrape"
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

type handoutData struct {
	Title    string
	Problems []*scrape.Problem
	Author   string
}

//go:embed templates/handout.html
var handout_template string

func handoutHandler(w http.ResponseWriter, r *http.Request) {
	// temporary
	tmpl, err := template.New("single").Parse(handout_template)
	if err != nil {
		panic(err)
	}

	ids := r.URL.Query()["id"]
	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")
	if len(title) == 0 {
		title = "An Anonymous Handout"
	}
	problems := make([]*scrape.Problem, len(ids))
	w.Header().Add("Content-Type", "text/html")
	for index, id := range ids {
		res, err := client.GetById(id)
		if err != nil {
			problems[index] = &scrape.Problem{
				Source:   "404 Not Found",
				Rendered: template.HTMLEscapeString(title),
			}
		} else {
			err := json.Unmarshal([]byte(res), &problems[index])
			if err != nil {
				panic(err)
			}
		}
	}
	err = tmpl.Execute(w, &handoutData{
		Title:    template.HTMLEscapeString(title),
		Author:   template.HTMLEscapeString(author),
		Problems: problems,
	})
	if err != nil {
		panic(err)
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
	mux.HandleFunc("/handout", handoutHandler)
	return mux
}
