package scrape

import (
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stderr, "[Scraper Info]  ", 0)

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
	w := sync.WaitGroup{}
	channel := make(chan []Problem, len(categories))
	for _, id := range categories {
		w.Add(1)
		go func(w *sync.WaitGroup, channel chan []Problem, id int) {
			resp, err := session.GetCategory(id)
			if err != nil {
				logger.Println("Error", err)
			} else {
				r := resp.ToProblems(session)
				channel <- r
			}
			w.Done()
		}(&w, channel, id)
	}
	w.Wait()
	res := make([]Problem, 0)
	close(channel)
	for c := range channel {
		res = append(res, c...)
	}
	return res
}
