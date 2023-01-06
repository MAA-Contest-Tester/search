package scrape

type Problem struct {
	Url        string `json:"url"`
	Source     string `json:"source"`
	Statement  string `json:"statement"`
	Solution   string `json:"solution"`
	Categories string `json:"categories"`
}
