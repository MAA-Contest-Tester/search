package database

import (
	"encoding/json"
	"log"
	"os"
	"regexp"

	"github.com/MAA-Contest-Tester/search/scrape"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/google/uuid"
)

var logger = log.New(os.Stderr, "[Redis Client]  ", 0);

type SearchClient struct {
  client *redisearch.Client
}

func Client() *SearchClient {
  db_host := "";
  if host,exists := os.LookupEnv("REDIS"); exists {
    db_host = host;
  } else {
    db_host = "localhost:6379";
  }
  client := redisearch.NewClient(db_host, "Problems");
  return &SearchClient{client: client};
}

func (c *SearchClient) AddProblems(problems []scrape.Problem) {
  schema := redisearch.NewSchema(redisearch.DefaultOptions).
            AddField(redisearch.NewTextField("url")).
            AddField(redisearch.NewTextField("statement")).
            AddField(redisearch.NewTextField("solution")).
            AddField(redisearch.NewTextField("source"))
  c.client.Drop();
  if err := c.client.CreateIndex(schema); err != nil {
    logger.Fatal(err);
  }

  docs := make([]redisearch.Document, 0);
  whitespace := regexp.MustCompile(`^\s*$`);
  for _,p := range problems {
    // do not include problems in the database that contain literally nothing.
    if whitespace.Match([]byte(p.Statement)) || whitespace.Match([]byte(p.Solution)) {
      continue
    }
    doc := redisearch.NewDocument(uuid.NewString(), problemScore(p));
    doc.Set("url", p.Url).
        Set("statement", p.Statement).
        Set("solution", p.Solution).
        Set("source", p.Source)
    docs = append(docs, doc);
  }
  if err := c.client.Index(docs...); err != nil {
		logger.Println("From Document Insertion ", err)
	}
}

func (c *SearchClient) Search(query string) (string, error) {
  q := redisearch.NewQuery(query);
  docs, _, error := c.client.Search(q.Limit(0,15))
  if error != nil {
    return "[]", error;
  }
  res := []map[string]interface{}{};
  for _,d := range docs {
    res = append(res, d.Properties);
  }
  out,err := json.Marshal(res);
  if err != nil {
    log.Fatal(err);
  }
  return string(out), nil;
}
