package database

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/MAA-Contest-Tester/search/backend/scrape"
	"github.com/meilisearch/meilisearch-go"
)

type MeiliSearchClient struct {
	client *meilisearch.Client
	index  *meilisearch.Index
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
	index := client.Index("problems")
	return &MeiliSearchClient{client: client, index: index}
}

func (c *MeiliSearchClient) Drop() {
	c.index.DeleteAllDocuments()
}

func calculateSynonyms() map[string][]string {
	synonyms := map[string][]string{
		"isl":           {"imo shortlist"},
		"imo shortlist": {"isl"},
		"sl":            {"shortlist"},
		"shortlist":     {"sl"},

		"bmo":              {"balkan mo"},
		"rmm":              {"romanian masters"},
		"hmmt":             {"Harvard MIT Mathematics"},
		"smt":              {"Stanford Mathematics Tournament"},
		"pumac":            {"Princeton University Math"},
		"jmo":              {"usajmo"},
		"amo":              {"usamo"},
		"mpfg":             {"math prize girls"},
		"math prize girls": {"mpfg"},

		"geo":   {"geometry"},
		"alg":   {"algebra"},
		"nt":    {"number theory"},
		"combo": {"combinatorics"},
		"fe":    {"functional equation"},
	}
	// A1 => Algebra 1, G8 => Geometry 8, ...
	categories := map[string]string{"a": "algebra", "g": "geometry", "n": "number theory", "c": "combinatorics"}
	for key, value := range categories {
		for i := 1; i < 12; i++ {
			short := fmt.Sprintf("%v%v", key, i)
			long := fmt.Sprintf("%v %v", value, i)
			synonyms[short] = []string{long}
			synonyms[long] = []string{short}
		}
	}
	return synonyms
}

func (c *MeiliSearchClient) AddProblems(problems []scrape.Problem) {
	docs := make([]map[string]interface{}, 0)
	order := []string{"source", "categories", "statement"}
	settings := meilisearch.Settings{
		SearchableAttributes: []string{"source", "categories", "statement"},
		Synonyms:             calculateSynonyms(),
		RankingRules: []string{
			"attribute",
			"words",
			"proximity",
			"typo",
			"sort",
			"exactness",
			"year:desc",
		},
		StopWords: []string{
			"a", "is", "the", "an", "and", "as", "at", "be", "but", "by",
			"into", "it", "not", "of", "on", "or", "their",
			"there", "these", "they", "this", "to", "was", "will",
		},
	}
	c.index.UpdateSettings(&settings)
	_, err := c.index.UpdateSearchableAttributes(&order)
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
		h := sha256.New()
		h.Write([]byte(doc["rendered"].(string)))

		doc["id"] = fmt.Sprintf("%x",h.Sum(nil))
		doc["year"] = extractYear(p)

		docs = append(docs, doc)
	}
	task, err := c.index.AddDocuments(docs)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Println("Task", task.TaskUID)
}

func (c *MeiliSearchClient) Search(query string) (string, error) {
	result, err := c.index.Search(query, &meilisearch.SearchRequest{
		Limit: 20,
	})
	if err != nil {
		return "[]", err
	}
	out, _ := json.Marshal(result.Hits)
	return string(out), nil
}
