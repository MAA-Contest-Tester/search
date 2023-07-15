package scrape

type Problem struct {
	Url        string `json:"url"`
	Source     string `json:"source"`
	Statement  string `json:"statement"`
	Rendered  string `json:"rendered"`
	Solution   string `json:"solution"`
	Categories string `json:"categories"`
}

type Contest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// key is general category (e.g. olympiad, USA contests, etc)
type ContestList map[string][]Contest

type Meta struct {
	Contests     ContestList `json:"contestlist"`
	ProblemCount int         `json:"problemcount"`
	Date         string      `json:"date"`
}

type ScrapeResult struct {
	Meta     Meta      `json:"meta"`
	Problems []Problem `json:"problems"`
}
