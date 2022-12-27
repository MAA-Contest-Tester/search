package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Addition struct {
	search *regexp.Regexp
	insert string
}

var additions = []Addition{
	{search: regexp.MustCompile("introductory"), insert: "beginner easy"},
	{search: regexp.MustCompile("intermediate"), insert: "medium middle"},
	{search: regexp.MustCompile("olympiad"), insert: "proof hard difficult"},
	{search: regexp.MustCompile("geometry"), insert: "geo"},
	{search: regexp.MustCompile("combinatorics"), insert: "counting combo"},
	{search: regexp.MustCompile("number theory"), insert: "nt mod"},
	{search: regexp.MustCompile("inequality"), insert: "algebra bound"},
	{search: regexp.MustCompile("trigonometry"), insert: "trig algebra geometry"},
}

func modifyWithAdditions(s string) string {
	s = strings.ToLower(s)
	res := []string{s}
	for _, a := range additions {
		if a.search.Match([]byte(s)) {
			res = append(res, a.insert)
		}
	}
	return strings.Join(res, " ")
}

var splitSolutionURL = regexp.MustCompile(`index\.php|\?`)

func Categorize(solution_url string) string {
	if len(solution_url) == 0 || redlink.Match([]byte(solution_url)) {
		return ""
	}
	type APIJson struct {
		Parse map[string]interface{} `json:"parse"`
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
	resp, err := http.Get(url)
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
	categories := data.Parse["categories"]
	res := make([]string, 0)
	for _, c := range categories.([]interface{}) {
		category := strings.ReplaceAll(c.(map[string]interface{})["*"].(string), "_", " ")
		res = append(res, modifyWithAdditions(category))
	}
	return strings.Join(res, " ")
}
