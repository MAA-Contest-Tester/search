package scrape

import (
	"log"
	"sync"
)

/*
[]int {
	
	3416,   // aime
	3414,   // amc 10
	3415,   // amc 12
	3413,   // amc 8
	3427,   // mpfg
	953466, // mpfg olympiad

	3412,   // usamts
	3409,   // usamo
	3420,   // usajmo
	3429,   // elmo
	3222,   // imo
	3227,   // jbmo
	3411,   // usa tst
	3424,   // usa tstst
	3282,   // china tst
	3223,   // imo shortlist
	3226,   // apmo
	3246,   // egmo
	3225,   // balkan mo
	3372,   // sharygin
	3238,   // rmm
	3277,   // canada mo
	3383,   // kmo
	603052, // kjmo
	3284,   // china mo
	3287,   // cgmo
	3288,   // china second round
	3384,   // korea final round

	915845, //balkan mo shortlist
	3231,   //baltic way
	3371,   //all-russian olympiad

	2746308, // chmmc
	253928,  // cmimc
	3417,    // hmmt,
	2881068, // hmmt november
	3418,    // smt,
	2503467, // bmt
	3426,    // pumac
	233906,  // bamo
}
*/

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
	return ScrapeResult {
		Contests: contestlist,
		Problems: session.ScrapeForumList(res),
	}
}
