package scrape

import (
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stderr, "[Scraper Info]  ", 0);

func scrapeSingle(url string, ch chan []Problem, wg *sync.WaitGroup) {
  logger.Println("Scraping", url , "...");
  ch <- ScrapeAops(url);
  wg.Done();
  logger.Println("Done Scraping", url);
}

func ScrapeList(problemsets []string) []Problem {
  w := sync.WaitGroup{};
  channel := make(chan []Problem, len(problemsets));
  for _,url := range problemsets {
    w.Add(1);
    go scrapeSingle(url, channel, &w);
  }
  w.Wait();
  close(channel);
  res := make([]Problem, 0);
  for c := range channel {
    res = append(res, c...);
  }
  return res;
}
