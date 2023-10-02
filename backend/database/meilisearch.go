package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/MAA-Contest-Tester/search/backend/scrape"
	"github.com/meilisearch/meilisearch-go"
)

type MeiliSearchClient struct {
	client        *meilisearch.Client
	problemsIndex *meilisearch.Index
	contestsIndex *meilisearch.Index
}

func InitMeiliSearchClient() *MeiliSearchClient {
	host := "http://localhost:7700"
	if redis_host, exists := os.LookupEnv("MEILISEARCH"); exists {
		host = redis_host
	}
	key := ""
	if meilisearch_key, exists := os.LookupEnv("MEILISEARCH_KEY"); exists {
		key = meilisearch_key
	}
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   host,
		APIKey: key,
	})
	return &MeiliSearchClient{
		client:        client,
		problemsIndex: client.Index("problems"),
		contestsIndex: client.Index("contests"),
	}
}

func (c *MeiliSearchClient) Drop() {
	c.problemsIndex.DeleteAllDocuments()
}

func calculateSynonyms() map[string][]string {
	synonyms := map[string][]string{
		"isl":           {"imo shortlist"},
		"imo shortlist": {"isl"},
		"sl":            {"shortlist"},
		"shortlist":     {"sl"},
		"mo":            {"math-olympiad", "national-olympiad"},

		"tst":                 {"team selection test"},
		"team selection test": {"tst"},

		"bmo":              {"balkan mo"},
		"bamo":             {"bay-area-mathematical-olympiad"},
		"rmm":              {"romanian masters"},
		"hmmt":             {"Harvard-MIT-Mathematics-Tournament"},
		"hmnt":             {"Harvard-MIT-Mathematics-Tournament-November"},
		"smt":              {"stanford-mathematics-tournament"},
		"bmt":              {"berkeley-math-tournament"},
		"pumac":            {"princeton-university-math-competition"},
		"jmo":              {"usajmo"},
		"amo":              {"usamo"},
		"mpfg":             {"math-prize-for-girls"},
		"math prize girls": {"mpfg"},

		"geo":   {"geometry", "geometrical"},
		"alg":   {"algebra"},
		"nt":    {"number theory"},
		"combo": {"combinatorics", "combinatorial"},
		"fe":    {"functional equation"},
		"mmp":   {"method-of-moving-points"},
	}
	// A1 => Algebra 1, G8 => Geometry 8, ...
	categories := map[string]string{"a": "algebra", "g": "geometry", "n": "number-theory", "c": "combinatorics"}
	for key, value := range categories {
		for i := 1; i < 12; i++ {
			short := fmt.Sprintf("%v%v", key, i)
			long := fmt.Sprintf("%v-problem-%v", value, i)
			longer := fmt.Sprintf("%v-problem-%v%v", value, key, i)
			synonyms[short] = []string{long, longer}
			synonyms[long] = []string{short, longer}
			synonyms[longer] = []string{short, long}
		}
	}
	return synonyms
}

func (c *MeiliSearchClient) AddProblems(problems []scrape.Problem) {
	docs := make([]map[string]interface{}, 0)
	order := []string{"source", "categories", "statement"}
	settings := meilisearch.Settings{
		SearchableAttributes: []string{"source", "statement", "categories"},
		Synonyms:             calculateSynonyms(),
		RankingRules: []string{
			"words",
			"attribute",
			"exactness",
			"typo",
			"proximity",
			"sort",
			"year:desc",
		},
		StopWords: []string{
			"a", "is", "the", "an", "and", "as", "at", "be", "but", "by",
			"into", "it", "not", "of", "on", "or", "their",
			"there", "these", "they", "this", "to", "was", "will", "for",
			"problem",
		},
	}
	c.problemsIndex.UpdateSettings(&settings)
	_, err := c.problemsIndex.UpdateSearchableAttributes(&order)
	if err != nil {
		logger.Fatalln(err)
	}

	whitespace := regexp.MustCompile(`^\s*$`)
	for _, p := range problems {
		// do not include problems in the database that contain literally nothing.
		if whitespace.Match([]byte(p.Statement)) || whitespace.Match([]byte(p.Solution)) {
			continue
		}
		marshaled, err := json.Marshal(p)
		if err != nil {
			logger.Fatal(err)
		}
		doc := map[string]interface{}{}
		err = json.Unmarshal(marshaled, &doc)
		if err != nil {
			logger.Fatal(err)
		}
		// TODO: set doc[id] to a urlized version of the source.

		doc["id"] = SourceToId(doc["source"].(string))
		doc["year"] = ExtractYear(p)

		docs = append(docs, doc)
	}
	task, err := c.problemsIndex.AddDocuments(docs)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Println("Task", task.TaskUID)
}

func (c *MeiliSearchClient) SearchProblems(query string, offset int) (string, error) {
	result, err := c.problemsIndex.Search(query, &meilisearch.SearchRequest{
		Limit:                 20,
		AttributesToHighlight: []string{"source", "categories"},
		HighlightPreTag:       "<span class=\"highlight\">",
		HighlightPostTag:      "</span>",
		Offset:                int64(offset),
	})
	if err != nil {
		return "[]", err
	}
	out, _ := json.Marshal(result.Hits)
	return string(out), nil
}

func (c *MeiliSearchClient) GetById(id string) (string, error) {
	var p map[string]interface{}
	message := c.problemsIndex.GetDocument(id, &meilisearch.DocumentQuery{}, &p)
	if message != nil {
		return "{}", errors.New("404 Not Found")
	}
	text, _ := json.Marshal(&p)
	return string(text), nil
}
