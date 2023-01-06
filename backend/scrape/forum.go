package scrape

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

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
	//f := ForumClient{};
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
		Timeout: time.Minute * 40,
	}
	resp, err := client.Do(f.InitRequest(url.Values{"a": {"fetch_topic"}, "topic_id": {strconv.Itoa(id)}}))
	if err != nil || resp == nil || resp.StatusCode != 200 {
		logger.Println(err)
		return nil, err
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

func (resp *CategoryResponse) ToProblems(f *ForumSession) []Problem {
	type Topic struct {
		Problem Problem
		Id      int
	}
	items := resp.Response.Category.Items
	problems := make([]Topic, 0)
	front_label := ""
	// make sure we're not dealing with Solutions
	solution_re := regexp.MustCompile(`[Ss]olution`)
	if solution_re.Match([]byte(resp.Response.Category.Name)) {
		return []Problem{}
	}
	for _, p := range items {
		announcement := p.PostData.CategoryId == 75
		label := p.PostData.CategoryId == resp.Response.Category.CategoryId
		notpost := strings.ToLower(p.Type) != "post"
		if label {
			front_label = p.PostData.Rendered
		}
		if announcement || label || notpost {
			continue
		}
		problem := Problem{
			Url: fmt.Sprintf(
				"https://artofproblemsolving.com/community/c%v",
				resp.Response.Category.CategoryId,
			),
			Source: fmt.Sprintf(
				"%v %v Problem %v",
				resp.Response.Category.Name,
				front_label,
				p.Title,
			),
			Statement: p.PostData.Rendered,
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
			x.Problem.Categories = strings.Join(tags, " ")
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

func parseProblemRenderedHTML(text string) (string, error) {
	nd, err := html.Parse(strings.NewReader(text))
	if err != nil {
		return "", err
	}
	// replace with dollar signs
	for child := nd.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			child.Data = strings.ReplaceAll(child.Data, "$", `$\textdollar$`)
		}
	}
	doc := goquery.NewDocumentFromNode(nd)
	doc.Find("img[alt].latex, img[alt].latexcenter").Each(func(i int, s *goquery.Selection) {
		s.SetText(s.AttrOr("alt", ""))
	})
	doc.Find("img[src].asy-image").Each(func(i int, s *goquery.Selection) {
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
	t := doc.Text()
	return t, nil
}

func (f *ForumSession) GetCategory(id int) (*CategoryResponse, error) {
	logger.Println("Parsing Forum Category", id, "...")
	client := http.Client{
		Timeout: time.Minute * 20,
	}
	resp, err := client.Do(f.InitRequest(url.Values{"a": {"fetch_category_data"}, "category_id": {strconv.Itoa(id)}}))
	if err != nil {
		logger.Println(err)
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
			r, err := parseProblemRenderedHTML(x.PostData.Rendered)
			if err != nil {
				logger.Fatal(err)
			}
			serialized.Response.Category.Items[i].PostData.Rendered = r
		}
		logger.Println("Finished Parsing Forum Category", id)
		return &serialized, nil
	}
}
