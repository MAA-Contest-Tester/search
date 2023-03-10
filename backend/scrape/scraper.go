package scrape

import (
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stderr, "[Scraper Info]  ", 0)

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
