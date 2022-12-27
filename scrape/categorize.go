package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var splitSolutionURL = regexp.MustCompile(`index\.php|\?`)



func Categorize(solution_url string) string {
	if len(solution_url) == 0 || redlink.Match([]byte(solution_url)) {
		return "";
	}
	type APIJson struct {
		Parse map[string]interface{} `json:"parse"`
	}
	page := splitSolutionURL.Split(solution_url, -1)[1]
	if len(page) == 0 {
		logger.Println(solution_url, page);
		return "";
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
	defer resp.Body.Close();
	body, err := io.ReadAll(resp.Body);
	if err != nil {
		logger.Println("Err while reading body", url, err)
		return "";
	}
	data := APIJson{};
	err = json.Unmarshal(body, &data); if err != nil {
		logger.Println("Err while reading json", url, err)
	}
	categories := data.Parse["categories"]
	res := make([]string, 0);
	for _,c := range categories.([]interface{}) {
		res = append(res, strings.ReplaceAll(c.(map[string]interface{})["*"].(string), "_", " "));
	}
	return strings.Join(res, " ");
}
