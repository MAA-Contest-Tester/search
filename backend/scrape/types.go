package scrape

type Problem struct {
	Url        string `json:"url"`
	Source     string `json:"source"`
	Statement  string `json:"statement"`
	Solution   string `json:"solution"`
	Categories string `json:"categories"`
}

type Contest struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

// key is general category (e.g. olympiad, USA contests, etc)
type ContestList map[string][]Contest

type ScrapeResult struct {
	Contests ContestList
	Problems []Problem `json:"problems"`
}
