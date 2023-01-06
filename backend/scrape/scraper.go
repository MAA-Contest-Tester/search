package scrape

import (
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stderr, "[Scraper Info]  ", 0)
var WorkerCount int

func ScrapeWikiList(problemsets []string) []Problem {
	w := sync.WaitGroup{}
	channel := make(chan []Problem, len(problemsets))
	for _, url := range problemsets {
		w.Add(1)
		go func(url string, ch chan []Problem, wg *sync.WaitGroup) {
			logger.Println("Scraping", url, "...")
			ch <- ScrapeAops(url)
			wg.Done()
			logger.Println("Done Scraping", url)
		}(url, channel, &w)
	}
	w.Wait()
	close(channel)
	res := make([]Problem, 0)
	for c := range channel {
		res = append(res, c...)
	}
	return res
}

func (session *ForumSession) ScrapeForumList(categories []int) []Problem {
	logger.Println("workers", WorkerCount)
	jobs := make(chan int, len(categories))
	results := make(chan []Problem, len(categories))

	worker := func(jobs <-chan int, results chan<- []Problem) {
		for id := range jobs {
			resp, err := session.GetCategory(id)
			if err != nil {
				logger.Println("Error", err)
			} else {
				r := resp.ToProblems(session)
				results <- r
			}
		}
	}
	for i := 0; i < WorkerCount; i++ {
		go worker(jobs, results)
	}
	for _, id := range categories {
		jobs <- id
	}
	close(jobs)
	res := make([]Problem, 0)
	for i := 0; i < len(categories); i++ {
		res = append(res, (<-results)...)
	}
	return res
}
