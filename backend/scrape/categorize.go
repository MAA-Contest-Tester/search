package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Addition struct {
	search *regexp.Regexp
	insert string
}

func insertAdditions(s string, additions []Addition) string {
	s = strings.ToLower(s)
	res := []string{s}
	for _, a := range additions {
		if a.search.Match([]byte(s)) {
			res = append(res, a.insert)
		}
	}
	return strings.Join(res, " ")
}

var categoryAdditions = []Addition{
	{search: regexp.MustCompile("introductory"), insert: "beginner easy"},
	{search: regexp.MustCompile("intermediate"), insert: "medium middle"},
	{search: regexp.MustCompile("olympiad|mo|shortlist|isl"), insert: "proof hard difficult"},
	{search: regexp.MustCompile("geometry|g[1-9]"), insert: "geo"},
	{search: regexp.MustCompile("combinatorics|c[1-9]"), insert: "counting combo"},
	{search: regexp.MustCompile("number theory|n[1-9]|nt"), insert: "nt mod"},
	{search: regexp.MustCompile("algebra|a[1-9]"), insert: "algebra"},
	{search: regexp.MustCompile("inequality"), insert: "algebra bound"},
	{search: regexp.MustCompile("trigonometry"), insert: "trig algebra geometry"},
	{search: regexp.MustCompile("imo shortlist|isl"), insert: "imo shortlist isl"},
}

var statementAdditions = []Addition{
	{search: regexp.MustCompile(`triangle|square|quadrilateral|cyclic|circum|incenter|acute`), insert: "geometry geo"},
	{search: regexp.MustCompile(`gcd|prime|gcd|lcm|divisor`), insert: "number theory nt"},
	{search: regexp.MustCompile(`[a-z]\^[2-9]|sequence|function|polynomial|inequality`), insert: "algebra alg"},
	{search: regexp.MustCompile(`probability|choose|game|rows|columns`), insert: "algebra alg"},
}

var splitSolutionURL = regexp.MustCompile(`index\.php|\?`)

func CategorizeWiki(solution_url string) string {
	if len(solution_url) == 0 || redlink.Match([]byte(solution_url)) {
		return ""
	}
	type APIJson struct {
		Parse struct {
			Categories []struct {
				Name string `json:"*"`
			} `json:"categories"`
		} `json:"parse"`
	}
	page := splitSolutionURL.Split(solution_url, -1)[1]
	if len(page) == 0 {
		logger.Println(solution_url, page)
		return ""
	}
	if page[0] == '/' {
		page = page[1:]
	}
	url := fmt.Sprintf(
		"https://artofproblemsolving.com/wiki/api.php?action=parse&format=json&page=%v", page,
	)
	client := http.Client{
		Timeout: time.Minute * 20,
	}
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Println("Err while reading body", url, err)
		return ""
	}
	data := APIJson{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Println("Err while reading json", url, err)
	}
	res := make([]string, 0)
	for _, c := range data.Parse.Categories {
		category := strings.ReplaceAll(c.Name, "_", " ")
		res = append(res, insertAdditions(category, categoryAdditions))
	}
	return strings.Join(res, " ")
}

func CategorizeForum(problem *Problem) {
	problem.Categories = insertAdditions(problem.Source, categoryAdditions) + insertAdditions(problem.Statement, statementAdditions)
}
