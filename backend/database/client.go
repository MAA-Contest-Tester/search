package database

import (
	"log"
	"os"

	"github.com/MAA-Contest-Tester/search/backend/scrape"
)

var logger = log.New(os.Stderr, "[Database Client]  ", 0)

type SearchClient interface {
	Drop()
	AddProblems(problems []scrape.Problem)
	// returns either a json with all of the search results or an error.
	Search(query string) (string, error)
}
