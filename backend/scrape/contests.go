package scrape

import (
	"log"
	"sync"
	"time"
)

func (session *ForumSession) scrapeForumPage(id int) []int {
	resp, err := session.GetCategoryItems(id)
	if err != nil {
		log.Println("err", err)
		return []int{}
	}
	res := make([]int, 0)
	for _, x := range resp.Response.Category.Items {
		res = append(res, x.PostId)
	}
	return res
}

func ScrapeForumCategories(contestlist ContestList) ScrapeResult {
	session := InitForumSession()
	res := make([]int, 0)

	contestlist_length := 0
	for _, contests := range contestlist {
		contestlist_length += len(contests)
	}
	channel := make(chan []int, contestlist_length)
	wg := sync.WaitGroup{}
	for _, contests := range contestlist {
		for _, contest := range contests {
			wg.Add(1)
			go func(w *sync.WaitGroup, ch chan []int, id int) {
				ch <- session.scrapeForumPage(id)
				w.Done()
			}(&wg, channel, contest.Id)
		}
	}
	wg.Wait()
	close(channel)
	for x := range channel {
		res = append(res, x...)
	}
	problems := session.ScrapeForumList(res)
	return ScrapeResult{
		Meta: Meta{
			Contests:     contestlist,
			ProblemCount: len(problems),
			Date:         time.Now(),
		},
		Problems: session.ScrapeForumList(res),
	}
}

func (session *ForumSession) ScrapeForumList(contests []int) []Problem {
	w := sync.WaitGroup{}
	channel := make(chan []Problem, len(contests))
	for _, id := range contests {
		w.Add(1)
		go func(w *sync.WaitGroup, channel chan []Problem, id int) {
			resp, err := session.GetCategoryItems(id)
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
