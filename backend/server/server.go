package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	o := r.URL.Query().Get("offset")
	offset, _ := strconv.Atoi(o)
	w.Header().Add("Content-Type", "application/json")
	if q != "" {
		res, err := client.Search(q, offset)
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
	Author   string
	Description    string
	Problems []*scrape.Problem
}

//go:embed templates/handout.html
var handout_template string

func handoutHandler(w http.ResponseWriter, r *http.Request) {
	// temporary
	tmpl, err := template.New("single").Parse(handout_template)
	if err != nil {
		panic(err)
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "text/plain")
		io.WriteString(w, fmt.Sprintf("Invalid Request Method %v", r.Method))
		return
	}
	r.ParseForm()
	params := r.Form
	ids := params["id"]
	title := params.Get("title")
	author := params.Get("author")
	description := params.Get("description")
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
		Description:   template.HTMLEscapeString(description),
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
