package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"path"
	"sort"

	"github.com/MAA-Contest-Tester/search/backend/scrape"
)

const (
	pointsInTraining = 400
)

func fullQueries() [][]string {
	problems := scrape.ScrapeForumDefaults()
	data := [][]string{{"source", "statement"}}
	for _, p := range problems {
		if len(p.Statement) < 1000 {
			data = append(data, []string{p.Source, p.Statement})
		}
	}
	return data
}

func shuffleArray(ar [][]string) [][]string {
	type shuffleType struct {
		Element *[]string
		Index   int
		Value   float64
	}
	pairs := make([]shuffleType, len(ar))
	for i, x := range ar {
		rnd := rand.Float64()
		log.Println(rnd)
		pairs[i] = shuffleType{
			Element: &x,
			Index:   i,
			Value:   rnd,
		}
	}
	sort.Slice(pairs, func(i int, j int) bool {
		return pairs[i].Value < pairs[j].Value
	})
	for _, x := range pairs {
		(ar)[x.Index] = *(x.Element)
	}
	return ar
}

func getQueries(output string) {
	records := fullQueries()
	data := [][]string{{"source", "statement"}}
	perm := rand.Perm(len(records) - 1)
	for i := 0; i < pointsInTraining && i < len(perm); i++ {
		data = append(data, records[perm[i]+1])
	}
	log.Println(len(data))
	err := os.MkdirAll(path.Dir(output), 0755)
	if err != nil {
		log.Fatal("MkdirAll", err)
	}
	outputFile, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("OpenFile", err)
	}
	err = csv.NewWriter(outputFile).WriteAll(data)
	if err != nil {
		log.Fatal("WriteAll", err)
	}
}
