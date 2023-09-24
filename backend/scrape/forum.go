package scrape

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var logger = log.New(os.Stderr, "[Scraper Info]  ", 0)

type ForumSession struct {
	SessionId string `json:"id"`
	UserId    int    `json:"user_id"`
	Username  string `json:"username"`
	LoggedIn  bool   `json:"logged_in"`
	Role      string `json:"role"`
	Sid       string `json:",omitempty"`
}

func InitForumSession() ForumSession {
	sessionre := regexp.MustCompile(`AoPS\.session = ({.*?})`)
	resp, err := http.Get("https://artofproblemsolving.com")
	if err != nil {
		logger.Fatal(err)
	}
	data := ForumSession{}
	for _, c := range resp.Cookies() {
		if c.Name == "aopssid" {
			data.Sid = c.Value
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}
	session := sessionre.FindSubmatch(body)[1]
	if err = json.Unmarshal(session, &data); err != nil {
		logger.Fatal(err)
	}
	return data
}

func (f *ForumSession) InitRequest(body_input url.Values) *http.Request {
	body := url.Values{
		"aops_logged_in":  {strconv.FormatBool(f.LoggedIn)},
		"aops_user_id":    {strconv.Itoa(f.UserId)},
		"aops_session_id": {f.SessionId},
	}
	for k, v := range body_input {
		body[k] = v
	}
	req, err := http.NewRequest(
		http.MethodPost,
		"https://artofproblemsolving.com/m/community/ajax.php",
		strings.NewReader(body.Encode()),
	)
	if err != nil {
		logger.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "aopsuid", Value: strconv.Itoa(f.UserId)})
	req.AddCookie(&http.Cookie{Name: "aopssid", Value: f.Sid})
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	return req
}

type ErrorResponse struct {
	Code    string `json:"error_code,omitempty"`
	Message string `json:"error_msg,omitempty"`
}

/*
Parsing Topic Tags

E.g. https://artofproblemsolving.com/community/c6h1598717p9937285
*/

type TopicResponse struct {
	Response struct {
		Topic *struct {
			Tags []struct {
				Id   int    `json:"tag_id"`
				Text string `json:"tag_text"`
			} `json:"tags"`
		} `json:"topic"`
	} `json:"response"`
}

func (f *ForumSession) GetTopic(id int) (*TopicResponse, error) {
	logger.Println("Parsing Forum Topic", id, "...")
	client := http.Client{
		Timeout: time.Minute * 5,
	}
	resp, err := client.Do(f.InitRequest(url.Values{"a": {"fetch_topic"}, "topic_id": {strconv.Itoa(id)}}))
	if err != nil || resp == nil || resp.StatusCode != 200 {
		logger.Println(err)
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	x := ErrorResponse{}
	json.Unmarshal(respbody, &x)

	serialized := TopicResponse{}
	err = json.Unmarshal(respbody, &serialized)
	if err != nil || serialized.Response.Topic == nil {
		serializederror := ErrorResponse{}
		sererr := json.Unmarshal(respbody, &serializederror)
		if sererr != nil {
			log.Fatal(sererr)
		}
		if len(serializederror.Code) > 0 {
			return nil, errors.New(serializederror.Code)
		}
		return nil, err
	} else {
		logger.Println("Finished Parsing Forum Topic", id)
		return &serialized, nil
	}
}

/*
Parsing Problem Sets Per Category ID
E.g. https://artofproblemsolving.com/community/c3948_1997_imo_shortlist

Postdata will be
E.g. https://artofproblemsolving.com/community/c3223
*/

type Post struct {
	PostId   int    `json:"item_id"`
	Title    string `json:"item_text"`
	Type     string `json:"item_type"`
	PostData struct {
		TopicId    int    `json:"topic_id"`
		PostId     int    `json:"post_id"`
		CategoryId int    `json:"category_id"`
		Rendered   string `json:"post_rendered"`
		Canonical  string `json:"post_canonical"`
	} `json:"post_data"`
}

type CategoryResponse struct {
	Response struct {
		Category *struct {
			CategoryId int    `json:"category_id"`
			Name       string `json:"category_name"`
			Items      []Post `json:"items"`
		} `json:"category"`
	} `json:"response"`
}

// function to clean out some of the BS people perform on C&P titles
func ProcessProblemSource(s string) string {
	// get rid of any non-alphanumeric characters.
	filtered := make([]rune, 0)
	for _, c := range []rune(s) {
		if c == '-' || c == '/' {
			filtered = append(filtered, ' ')
		} else if '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == ' ' {
			filtered = append(filtered, c)
		}
	}
	s = string(filtered)
	// fix 2017 IMO ShortiIst
	shortiIstregex := regexp.MustCompile(`Short[iI][iI]st`)
	s = shortiIstregex.ReplaceAllString(s, "Shortlist")
	// get rid of redundant "Problems"
	problemsRegex := regexp.MustCompile(`\s*[Pp]roblems*\s*`)
	s = problemsRegex.ReplaceAllString(s, " ")
	islRegex := regexp.MustCompile(`ISL`)
	s = islRegex.ReplaceAllString(s, "IMO Shortlist")
	return string(filtered)
}

/*
Helper function that eliminates completely useless HTML tags in titles and whatnot.
*/
func RemoveHtmlBS(s string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		logger.Fatal(err)
	}
	return doc.Text()
}

// disqualify tags that represent contests because they pollute search results.
var contestRegex []regexp.Regexp = []regexp.Regexp{
	*regexp.MustCompile(`.*sl$`),
	*regexp.MustCompile(`.*mo$`),
	*regexp.MustCompile(`.*[ms]t$`),
	*regexp.MustCompile(`amc|aime`),
	*regexp.MustCompile(`\d{4}`),
}

func ProcessTags(s string) string {
	words := strings.Fields(s)
	processed := make([]string, 0)
	seen := map[string]int{}
	for _, word := range words {
		include := true
		word = strings.ToLower(word)

		for _, re := range contestRegex {
			if re.Match([]byte(word)) {
				include = false
			}
		}
		if _, exists := seen[word]; exists {
			include = false
		}

		if include {
			processed = append(processed, word)
			seen[word] = 1
		}
	}
	return strings.Join(processed, " ")
}

func (resp *CategoryResponse) ToProblems(f *ForumSession) []Problem {
	type Topic struct {
		Problem Problem
		Id      int
	}
	items := resp.Response.Category.Items
	problems := make([]Topic, 0)

	front_label := ""
	// there are instances where there are two or more labels stacked on top of
	// each other: such as one line containing "I" and annother line containing
	// "(insert date)" for a specific AIME.
	previous_label := false
	// make sure we're not dealing with Solutions
	solution_re := regexp.MustCompile(`[Ss]olution`)
	if solution_re.Match([]byte(resp.Response.Category.Name)) {
		return []Problem{}
	}
	for _, p := range items {
		// the "These problems are copyright of MAA" message
		announcement := p.PostData.CategoryId == 75
		// When one of the rows is just a label saying "this is day 2"
		label := p.PostData.CategoryId == resp.Response.Category.CategoryId
		// Straight-up when not a post
		notpost := strings.ToLower(p.Type) != "post"
		if label {
			// only take the first label that comes in. redundant afterwards.
			if !previous_label {
				front_label = p.PostData.Rendered
			}
			previous_label = true
		} else {
			previous_label = false
		}
		if announcement || label || notpost {
			continue
		}
		problem := Problem{
			Source: RemoveHtmlBS(fmt.Sprintf(
				"%v %v Problem %v",
				// e.g. "2023 USAMO"
				ProcessProblemSource(resp.Response.Category.Name),
				// e.g. "Day 2"
				front_label,
				p.Title,
			)),
			Statement: p.PostData.Canonical,
			Rendered:  p.PostData.Rendered,
			Url: fmt.Sprintf(
				"https://artofproblemsolving.com/community/c%v",
				resp.Response.Category.CategoryId,
			),
			Solution: fmt.Sprintf(
				"https://artofproblemsolving.com/community/c%vh%vp%v",
				resp.Response.Category.CategoryId,
				p.PostData.TopicId,
				p.PostData.PostId,
			),
		}
		problems = append(problems, Topic{
			Problem: problem,
			Id:      p.PostData.TopicId,
		})
	}
	channel := make(chan Problem, len(problems))
	wg := sync.WaitGroup{}

	// fetch tags per category
	for _, x := range problems {
		wg.Add(1)
		go func(c chan Problem, w *sync.WaitGroup, x Topic) {
			t, err := f.GetTopic(x.Id)
			if err != nil || t == nil {
				channel <- x.Problem
				logger.Println(err)
				wg.Done()
				return
			}
			tags := make([]string, 0)
			for _, tag := range t.Response.Topic.Tags {
				tags = append(tags, tag.Text)
			}
			x.Problem.Categories = ProcessTags(
				strings.Join(tags, " "),
			)
			channel <- x.Problem
			wg.Done()
		}(channel, &wg, x)
	}

	wg.Wait()
	close(channel)
	// put all problems from the channel.
	res := make([]Problem, 0)
	for p := range channel {
		res = append(res, p)
	}
	return res
}

// helper function to half the width.
func reduceWidth(attr string) string {
	value, err := strconv.Atoi(attr)
	if err != nil {
		return attr
	} else {
		return strconv.Itoa(value / 2)
	}
}

/*
This function takes each problem statement (rendered as HTML) on AoPS and performs several processing steps:

- Have all asymptote images thet e
- Replace all image nodes with \includegraphics{...} so that it can be rendered by KaTeX
- Remove any images that are supposed to render LaTeX snippets and replace them with plain text snippets (i.e. $expression...$)
*/
func parseProblemRenderedHTML(text string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return "", err
	}
	for _, node := range doc.Nodes {
		// each node is structured as <html><head></head><body>{our text}</body></html>
		// so we have to structure it as below.
		node = node.FirstChild.FirstChild.NextSibling
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				// sometimes people put an actual dollar sign instead of
				// textdollar and this screws up all of the rendering later on
				child.Data = strings.ReplaceAll(child.Data, "$", `$\textdollar$`)
			}
		}
	}
	doc.Find("img[alt].latex, img[alt].latexcenter").Each(func(i int, s *goquery.Selection) {
		s.SetText(s.AttrOr("alt", ""))
	})
	doc.Find("img[src].asy-image, img[src].bbcode_img").Each(func(i int, s *goquery.Selection) {
		s.SetText(
			fmt.Sprintf(
				"$\\includegraphics[width=%v, height=%v, totalheight=%v]{https:%v}$",
				reduceWidth(s.AttrOr("width", "")),
				reduceWidth(s.AttrOr("height", "")),
				reduceWidth(s.AttrOr("height", "")),
				s.AttrOr("src", ""),
			),
		)
	})
	doc.Find("img[src].bbcode_img").Each(func(i int, s *goquery.Selection) {
		s.SetText(
			fmt.Sprintf(
				"$\\includegraphics[height=%v, totalheight=%v]{%v}$",
				"7em", "7em",
				s.AttrOr("src", ""),
			),
		)
	})
	t := doc.Text()
	return t, nil
}

/*

This is one function to take care of two different cases (but are the same
problem because of the recursive structure of AoPS categories):

1. Parsing out the problems from a specific year of a specific contest.
   e.g. Parsing all the problems from
   https://artofproblemsolving.com/community/c3381519 (The 2023 IMO Problems
   Category)
2. Parsing out all of the contest years of a specific contest
   e.g. Parsing out
   https://artofproblemsolving.com/community/c3223_imo_shortlist (The Collection
   that contains all IMO Shortlist Collections from every year)

*/

func (f *ForumSession) GetCategoryItems(id int) (*CategoryResponse, error) {
	logger.Println("Parsing Forum Category", id, "...")
	client := http.Client{
		Timeout: time.Minute * 5,
	}
	resp, err := client.Do(f.InitRequest(url.Values{"a": {"fetch_category_data"}, "category_id": {strconv.Itoa(id)}}))
	if err != nil {
		logger.Println(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	serialized := CategoryResponse{}
	err = json.Unmarshal(respbody, &serialized)
	if err != nil || serialized.Response.Category == nil {
		serializederror := ErrorResponse{}
		sererr := json.Unmarshal(respbody, &serializederror)
		if sererr != nil {
			log.Fatal(sererr)
		}
		if len(serializederror.Code) > 0 {
			return nil, errors.New(serializederror.Code)
		}
		return nil, err
	} else {
		for i, x := range serialized.Response.Category.Items {
			// in the second case described above, this effectively does
			// nothing.
			r, err := parseProblemRenderedHTML(x.PostData.Rendered)
			if err != nil {
				logger.Fatal(err)
			}
			serialized.Response.Category.Items[i].PostData.Canonical = r
		}
		logger.Println("Finished Parsing Forum Category", id)
		return &serialized, nil
	}
}
