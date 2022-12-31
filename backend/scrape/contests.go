package scrape

import (
	"log"
	"regexp"
	"sync"

	"github.com/gocolly/colly"
)

var wikicontests = []struct{ Url string; Search *regexp.Regexp }{
	{Url: "AIME_Problems_and_Solutions", Search: regexp.MustCompile(`/index.php/\d{4}_AIME(_I{1,2})?$`)},
	{Url: "AMC_10_Problems_and_Solutions", Search: regexp.MustCompile(`/index.php/\d{4}(_[A-Z,a-z]*?)?_AMC_10[AB]?$`)},
	{Url: "AMC_12_Problems_and_Solutions", Search: regexp.MustCompile(`/index.php/\d{4}(_[A-Z,a-z]*?)?_AMC_12[AB]?$`)},
	{Url: "AHSME_Problems_and_Solutions", Search: regexp.MustCompile(`/index.php/\d{4}_AHSME$`)},
	{Url: "USAJMO_Problems_and_Solutions", Search: regexp.MustCompile(`/index.php/\d{4}_USAJMO$`)},
	{Url: "USAMO_Problems_and_Solutions", Search: regexp.MustCompile(`/index.php/\d{4}_USAMO$`)},
	{Url: "IMO_Problems_and_Solutions", Search: regexp.MustCompile(`/index.php/\d{4}_IMO$`)},
	{Url: "JBMO_Problems_and_Solutions,_with_authors", Search: regexp.MustCompile(`/index.php/\d{4}_JBMO$`)},
	{Url: "AMC_8_Problems_and_Solutions", Search: regexp.MustCompile(`/index.php/\d{4}_(AMC_8|AJHSME)$`)},
}

var forums = []int {
	3411, // usa tst
	3424, // usa tstst
	3282, // china tst
	3223, // imo shortlist
	3226, // apmo
	3246, // egmo
	3225, // balkan mo
}

var redlink = regexp.MustCompile(`redlink=1`)

func scrapeWikiPage(url string, re *regexp.Regexp, w *sync.WaitGroup, channel chan []string) {
	url = "https://artofproblemsolving.com/wiki/index.php/" + url
	logger.Println("Parsing", url, "For Problemsets...")
	c := colly.NewCollector()

	res := make([]string, 0)

	c.OnHTML("div.mw-parser-output a[href]", func(el *colly.HTMLElement) {
		// fill out res.
		href := el.Attr("href")
		if redlink.Match([]byte(href)) {
			return
		}
		match := re.Match([]byte(href))
		if match {
			res = append(res, el.Request.AbsoluteURL(href)+"_Problems")
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		logger.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, ", Error:", err)
	})
	c.Visit(url)
	c.Wait()
	channel <- res
	w.Done()
	logger.Println("Finished Parsing", url)
}

func ScrapeWikiDefaults() []Problem {
	res := make([]string, 0)
	channel := make(chan []string, len(wikicontests))
	wg := sync.WaitGroup{}
	for _, contest := range wikicontests {
		wg.Add(1)
		go scrapeWikiPage(contest.Url, contest.Search, &wg, channel)
	}
	wg.Wait()
	close(channel)
	for x := range channel {
		res = append(res, x...)
	}
	return ScrapeWikiList(res);
}

func (session *ForumSession) scrapeForumPage(id int) []int {
	resp, err := session.GetCategory(id);
	if err != nil {
		log.Println("err", err);
		return []int{};
	}
	res := make([]int, 0);
	for _, x := range resp.Response.Category.Items {
		res = append(res, x.PostId);
	}
	return res;
}

func ScrapeForumDefaults() []Problem {
	session := InitForumSession();
	res := make([]int, 0);
	channel := make(chan []int, len(forums))
	wg := sync.WaitGroup{}
	for _, id := range forums {
		wg.Add(1);
		go func(w *sync.WaitGroup, ch chan []int, id int) {
			ch <- session.scrapeForumPage(id);
			w.Done();
		}(&wg, channel, id)
	}
	wg.Wait();
	close(channel);
	for x := range channel {
		res = append(res, x...)
	}
	return session.ScrapeForumList(res);
}
