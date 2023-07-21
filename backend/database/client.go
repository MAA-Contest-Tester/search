package database

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/MAA-Contest-Tester/search/backend/scrape"
)

var logger = log.New(os.Stderr, "[Database Client]  ", 0)

type SearchClient interface {
	Drop()
	AddProblems(problems []scrape.Problem)
	// returns either a json with all of the search results or an error.
	Search(query string, offset int) (string, error)
	GetById(id string) (string, error)
}

// common functions used for both meilisearch and redis
func SourceToId(source string) string {
	source = strings.ToLower(source)
	runes := []rune(source)
	res := []rune{}
	for _, r := range runes {
		switch {
		case ('0' <= r && r <= '9') || 'a' <= r && r <= 'z':
			res = append(res, r)
		case r == ' ':
			res = append(res, rune('-'))
		}
	}
	return string(res)
}

func ExtractYear(p scrape.Problem) float64 {
	re := regexp.MustCompile(`\d{4}`)
	num := re.Find([]byte(p.Source))
	if num != nil {
		res, _ := strconv.ParseFloat(string(num), 64)
		return res
	} else {
		return 1900.0
	}
}
