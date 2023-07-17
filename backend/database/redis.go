package database

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"regexp"

	"github.com/MAA-Contest-Tester/search/backend/scrape"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/google/uuid"
)

func problemScore(p scrape.Problem) float32 {
	// contests before this year should be penalized
	THRESHOLD := 2010.0
	return float32(1.0 / (1.0 + 0.1*math.Exp(THRESHOLD-ExtractYear(p))))
}

type RedisClient struct {
	client *redisearch.Client
}

func InitRedisClient() *RedisClient {
	host := "localhost:6379"
	if redis_host, exists := os.LookupEnv("REDIS"); exists {
		host = redis_host
	}
	client := redisearch.NewClient(host, "Problems")
	return &RedisClient{client: client}
}

func (c *RedisClient) Drop() {
	logger.Println("Dropping Database...")
	c.client.Drop()
}

func (c *RedisClient) AddProblems(problems []scrape.Problem) {
	options := redisearch.DefaultOptions.SetStopWords([]string{
		"a", "is", "the", "an", "and", "as", "at", "be", "but", "by",
		"into", "it", "not", "of", "on", "or", "their",
		"there", "these", "they", "this", "to", "was", "will",
	})
	schema := redisearch.NewSchema(*options).
		AddField(redisearch.NewTextField("url")).
		AddField(redisearch.NewTextField("statement")).
		AddField(redisearch.NewTextField("solution")).
		AddField(redisearch.NewTextField("source")).
		AddField(redisearch.NewTextField("categories"))

	if err := c.client.CreateIndex(schema); err != nil {
		logger.Println(err)
	}

	docs := make([]redisearch.Document, 0)
	whitespace := regexp.MustCompile(`^\s*$`)
	for _, p := range problems {
		// do not include problems in the database that contain literally nothing.
		if whitespace.Match([]byte(p.Statement)) || whitespace.Match([]byte(p.Solution)) {
			continue
		}
		doc := redisearch.NewDocument(uuid.NewString(), problemScore(p))
		doc.Set("url", p.Url).
			Set("statement", p.Statement).
			Set("statement", p.Statement).
			Set("solution", p.Solution).
			Set("source", p.Source).
			Set("categories", p.Categories)
		docs = append(docs, doc)
	}
	if err := c.client.Index(docs...); err != nil {
		logger.Println("From Document Insertion ", err)
	}
}

func (c *RedisClient) Search(query string) (string, error) {
	q := redisearch.NewQuery(query)
	docs, _, error := c.client.Search(q.Limit(0, 20))
	if error != nil {
		return "[]", error
	}
	res := []map[string]interface{}{}
	for _, d := range docs {
		res = append(res, d.Properties)
	}
	out, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	return string(out), nil
}

func (c *RedisClient) GetById(id string) (string, error) {
	p, err := c.client.Get(id)
	if err != nil {
		return "{}", err
	}
	serialized, err := json.Marshal(p)
	if err != nil {
		return "{}", err
	}
	return string(serialized), nil
}
