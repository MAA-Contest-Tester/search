package scrape

import (
	"regexp"
  "sync"
	"github.com/gocolly/colly"
);

type source struct {
  Url string
  Search *regexp.Regexp
}

var contests = []source{
  {Url: "AIME_Problems_and_Solutions",                        Search: regexp.MustCompile(`/index.php/\d{4}_AIME(_I{1,2})?$`)    },
  {Url: "AMC_10_Problems_and_Solutions",                      Search: regexp.MustCompile(`/index.php/\d{4}(_[A-Z,a-z]*?)?_AMC_10[AB]?$`)  },
  {Url: "AMC_12_Problems_and_Solutions",                      Search: regexp.MustCompile(`/index.php/\d{4}(_[A-Z,a-z]*?)?_AMC_12[AB]?$`)  },
  {Url: "AHSME_Problems_and_Solutions",                       Search: regexp.MustCompile(`/index.php/\d{4}_AHSME$`)   },
  {Url: "USAJMO_Problems_and_Solutions",                      Search: regexp.MustCompile(`/index.php/\d{4}_USAJMO$`)  },
  {Url: "USAMO_Problems_and_Solutions",                       Search: regexp.MustCompile(`/index.php/\d{4}_USAMO$`)   },
  {Url: "IMO_Problems_and_Solutions",                         Search: regexp.MustCompile(`/index.php/\d{4}_IMO$`)     },
  {Url: "JBMO_Problems_and_Solutions,_with_authors",          Search: regexp.MustCompile(`/index.php/\d{4}_JBMO$`)    },
  {Url: "AMC_8_Problems_and_Solutions",                       Search: regexp.MustCompile(`/index.php/\d{4}_(AMC_8|AJHSME)$`)  },
};

var redlink = regexp.MustCompile(`redlink=1`)

func scrapeContestPage(url string, re *regexp.Regexp, w *sync.WaitGroup, channel chan []string) {
  url = "https://artofproblemsolving.com/wiki/index.php/" + url;
  logger.Println("Parsing", url, "For Problemsets...")
  c := colly.NewCollector();
  
  res := make([]string, 0);

  c.OnHTML("div.mw-parser-output a[href]", func(el *colly.HTMLElement) {
    // fill out res.
    href := el.Attr("href")
    if redlink.Match([]byte(href)) {
      return;
    }
    match := re.Match([]byte(href));
    if match {
      res = append(res, el.Request.AbsoluteURL(href) + "_Problems");
    }
  });
  c.OnError(func(r *colly.Response, err error) {
		logger.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, ", Error:", err)
	})
  c.Visit(url);
  c.Wait();
  channel<-res;
  w.Done();
  logger.Println("Finished Parsing", url)
}

func ScrapeContestDefaults() []string {
  res := make([]string, 0);
  channel := make(chan []string, len(contests));
  wg := sync.WaitGroup{};
  for _,contest := range contests {
    wg.Add(1);
    go scrapeContestPage(contest.Url, contest.Search, &wg, channel);
  }
  wg.Wait();
  close(channel);
  for x := range channel {
    res = append(res,x...);
  }
  return res;
}
