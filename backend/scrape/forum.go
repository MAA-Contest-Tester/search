package scrape

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type ForumSession struct {
	SessionId string `json:"id"`
	UserId int `json:"user_id"`
	Username string `json:"username"`
	LoggedIn bool `json:"logged_in"`
	Role string `json:"role"`
	Sid string `json:",omitempty"`
}

func InitForumSession() ForumSession {
	sessionre := regexp.MustCompile(`AoPS\.session = ({.*?})`)
	//f := ForumClient{};
	resp, err := http.Get("https://artofproblemsolving.com")
	data := ForumSession{};
	for _,c := range resp.Cookies() {
		if c.Name == "aopssid" {
			data.Sid = c.Value;
		}
	}
	body, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		logger.Fatal(err);
	}
	session := sessionre.FindSubmatch(body)[1];
	if err = json.Unmarshal(session, &data); err != nil {
		logger.Fatal(err);
	}
	return data;
}

func (f *ForumSession) InitRequest(id int) *http.Request {
	body := url.Values {
		"category_id": { strconv.Itoa(id) },
		"aops_logged_in": { strconv.FormatBool(f.LoggedIn) },
		"a": { "fetch_category_data" },
		"aops_user_id": { strconv.Itoa(f.UserId) },
		"aops_session_id": { f.SessionId },
	}
	req,err := http.NewRequest(
		http.MethodPost,
		"https://artofproblemsolving.com/m/community/ajax.php",
		strings.NewReader(body.Encode()),
	);
	if err != nil {
		logger.Fatal(err);
	}
	req.AddCookie(&http.Cookie{Name: "aopsuid", Value: strconv.Itoa(f.UserId)})
	req.AddCookie(&http.Cookie{Name: "aopssid", Value: f.Sid})
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	return req;
}

type ErrorResponse struct {
	Code string `json:"errorcode,omitempty"`
	Message string `json:"error_msg,omitempty"`
}

/*
Parsing Problem Sets Per Category ID
E.g. https://artofproblemsolving.com/community/c3948_1997_imo_shortlist

Postdata will be 
E.g. https://artofproblemsolving.com/community/c3223
*/

type Post struct {
	PostId int `json:"item_id"`
	Title string `json:"item_text"`
	Type string `json:"item_type"`
	PostData struct {
		TopicId int `json:"topic_id"`
		PostId int `json:"post_id"`
		CategoryId int `json:"category_id"`
		Rendered string `json:"post_rendered"`
	} `json:"post_data"`
}

type CategoryResponse struct {
	Response struct {
		Category struct {
			CategoryId int `json:"category_id"`
			Name string `json:"category_name"`
			Items []Post `json:"items"`
		} `json:"category"`
	} `json:"response"`
}

func (resp *CategoryResponse) ToProblems() []Problem {
	items := resp.Response.Category.Items;
	res := make([]Problem,0)
	front_label := "";
	for _,p := range items {
		announcement := p.PostData.CategoryId == 75
		label := p.PostData.CategoryId == resp.Response.Category.CategoryId
		notpost := strings.ToLower(p.Type) != "post";
		if label {
			front_label = p.PostData.Rendered;
		}
		if announcement || label || notpost {
			continue
		}
		res = append(res, Problem{
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
			Categories: modifyWithAdditions(front_label),
		});
	}
	return res;
}

func parseProblemRenderedHTML(text string) (string, error) {
	nd, err := html.Parse(strings.NewReader(text));
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
		s.SetText(fmt.Sprintf("%v", s.AttrOr("src", "")[2:]))
	})
	t := doc.Text();
	return t, nil;
}

func (f *ForumSession) GetCategory(id int) (*CategoryResponse, error) {
	logger.Println("Parsing Forum Category", id, "...");
	client := http.Client{
		Timeout: time.Minute * 20,
	}
	resp, err := client.Do(f.InitRequest(id));
	if err != nil {
		logger.Println(err);
	}
	respbody, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		logger.Println(err);
		return nil, err;
	}
	serialized := CategoryResponse{};
	err = json.Unmarshal(respbody, &serialized); if err != nil {
		serializederror := ErrorResponse{};
		sererr := json.Unmarshal(respbody, &serializederror);
		if sererr != nil {
			log.Fatal(sererr);
		}
		if len(serializederror.Message) > 0 {
			return nil, errors.New(serializederror.Message);
		}
		return nil, err;
	} else {
		for i,x := range serialized.Response.Category.Items {
			r, err := parseProblemRenderedHTML(x.PostData.Rendered);
			if err != nil {
				logger.Fatal(err);
			}
			serialized.Response.Category.Items[i].PostData.Rendered = r;
		}
		logger.Println("Finished Parsing Forum Category", id);
		return &serialized, nil;
	}
}
