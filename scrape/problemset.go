package scrape

import (
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type Problem struct {
	Url       string `json:"url"`
	Source    string `json:"source"`
	Statement string `json:"statement"`
	Solution  string `json:"solution"`
	Categories  string `json:"categories"`
}

var statement_tags = []string{"p", "ol", "ul", "center"}

func isStatementTag(h *colly.HTMLElement) bool {
	s := h.Name
	for _, x := range statement_tags {
		if s == x {
			return true
		}
	}
	return false
}

func genSelector(s []string) string {
	for i := 0; i < len(s); i++ {
		s[i] = "div.mw-parser-output>" + s[i]
	}
	return strings.Join(s, ", ")
}

func makeSource(url string) string {
	re := regexp.MustCompile("[_#]")
	a := strings.Split(url, "/")
	last := a[len(a)-1]
	return re.ReplaceAllString(last, " ")
}

func ScrapeAops(url string) []Problem {
	c := colly.NewCollector()
	c.SetRequestTimeout(time.Minute * 10)
	latex_replace := func(_ int, b *colly.HTMLElement) {
		alt := b.Attr("alt")
		b.DOM.SetText(alt)
	}

	httperror := false
	res := make([]Problem, 0)

	selector := genSelector([]string{"h2", "h3", "p", "ol", "ul", "center"})
	c.OnHTML("html", func(el *colly.HTMLElement) {
		if httperror {
			logger.Println("Aborting Scraping because there is an HTTP Error.")
			return
		}
		// replace all images with latex form.
		el.ForEach(selector, func(idx int, content *colly.HTMLElement) {
			content.ForEach("img[alt]", latex_replace)
		})
		// fill out res.
		use_paragraph := false
		el.ForEach(selector, func(idx int, content *colly.HTMLElement) {
			if !isStatementTag(content) {
				span_id := content.ChildAttr("span[id]", "id")
				match, err := regexp.Match("Problem", []byte(span_id))
				if err != nil {
					logger.Printf("Regex Error, %v", err)
					os.Exit(1)
				}
				if match {
					res = append(res, Problem{
						Url:    content.Request.AbsoluteURL(url) + "#" + span_id,
						Source: makeSource(url + "#" + span_id),
					})
					use_paragraph = true
				} else {
					// if heading id doesn't start with "Problem", the text right after
					// isn't a problem statement so we should disregard it.
					use_paragraph = false
				}
			} else if use_paragraph {
				// terminates if there's "Solution" after a problem statement.
				match, err := regexp.Match(`^\s*[Ss]olution`, []byte(content.Text))
				if err != nil {
					logger.Printf("Regex Error, %v", err)
					os.Exit(1)
				}
				if match {
					res[len(res)-1].Solution = content.Request.AbsoluteURL(content.ChildAttr("a[href]", "href"))
					use_paragraph = false
				} else {
					res[len(res)-1].Statement = res[len(res)-1].Statement + content.Text
				}
			}
		})
	})
	c.OnError(func(r *colly.Response, err error) {
		logger.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, " Error:", err)
		httperror = true
	})
	c.Visit(url)
	c.Wait()

	// remove asymptote code
	asy_remove_regex := regexp.MustCompile(`\[asy\].*?\[/asy\]`)

	for i := 0; i < len(res); i++ {
		res[i].Statement = asy_remove_regex.ReplaceAllString(res[i].Statement, "")
		res[i].Statement = strings.TrimSpace(res[i].Statement)
	}

	// make categorization concurrent.
	type categoryResult struct {
		int
		string
	}
	categories := make(chan categoryResult, len(res));
	w := sync.WaitGroup{};
	for i := 0; i < len(res); i++ {
		w.Add(1);
		go func(i int, url string, c chan categoryResult, wg *sync.WaitGroup) {
			c <- categoryResult{i, Categorize(url)};
			wg.Done();
		}(i, res[i].Solution, categories, &w)
	}
	w.Wait();
	close(categories)
	i := 0;
	for category := range categories {
		res[category.int].Categories = category.string;
		i++;
	}
	return res
}
