package scrape

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

var logger = log.New(os.Stderr, "[Scraper Info]  ", 0)

// Cap concurrent AJAX calls to AoPS. The calls are issued via in-page fetch()
// inside a single Chromium instance.
const maxConcurrentRequests = 4

var requestSem = make(chan struct{}, maxConcurrentRequests)

// Hard rate cap. Empirical observation from local runs: Cloudflare lets the
// in-page fetch through for ~1000 successful requests, then walls the
// connection for ~20 seconds. Pacing at 5 req/sec keeps us under that budget
// indefinitely (300/min).
const minTimeBetweenRequests = 200 * time.Millisecond

var (
	rateMu        sync.Mutex
	nextRequestAt time.Time
)

func reserveRequestSlot() {
	rateMu.Lock()
	now := time.Now()
	target := nextRequestAt
	if now.After(target) {
		target = now
	}
	nextRequestAt = target.Add(minTimeBetweenRequests)
	delay := target.Sub(now)
	rateMu.Unlock()
	if delay > 0 {
		time.Sleep(delay)
	}
}

// Cloudflare may re-issue the JS challenge mid-scrape. Strategy:
//   - postAjax holds resolveMu.RLock for the duration of each fetch — many
//     fetches can run concurrently.
//   - When a fetch comes back with "Just a moment...", that goroutine calls
//     resolveChallenge which takes resolveMu.Lock — blocking all in-flight and
//     new fetches until the navigation is done.
//   - resolveChallenge dedupes: if another goroutine already resolved within
//     the last 10s, skip (the cookie is already fresh).
var (
	resolveMu       sync.RWMutex
	lastResolveAt   time.Time
	lastResolveAtMu sync.Mutex
)

// AoPS sits behind a Cloudflare JS challenge on /community/* that pins the
// resulting cf_clearance cookie to the solver's TLS fingerprint. The robust
// workaround is to do the bulk scrape from inside the same headless Chromium
// that solved the challenge — fetch() runs in the page context, so it inherits
// the cookie and the matching handshake automatically.
const browserUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

type ForumSession struct {
	SessionId string `json:"id"`
	UserId    int    `json:"user_id"`
	Username  string `json:"username"`
	LoggedIn  bool   `json:"logged_in"`
	Role      string `json:"role"`
	Sid       string `json:",omitempty"`

	browserCtx    context.Context
	cancelBrowser context.CancelFunc
	cancelAlloc   context.CancelFunc
}

func InitForumSession() ForumSession {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", "new"),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent(browserUserAgent),
	)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	bctx, cancelBrowser := chromedp.NewContext(allocCtx)

	logger.Println("Launching headless Chromium and solving Cloudflare challenge...")
	// Warm up the browser on bctx itself — chromedp lazy-launches Chrome on the
	// first Run() call and ties its lifetime to that call's context. If we used
	// a timeout child here, cancelling the timeout would also kill Chrome.
	if err := chromedp.Run(bctx); err != nil {
		cancelBrowser()
		cancelAlloc()
		logger.Fatal("chromedp launch: ", err)
	}

	// Use bctx (no timeout child) so the navigation state stays alive for
	// subsequent postAjax calls. chromedp ties some page state to the first
	// navigation's context — cancelling that context resets the tab.
	var sessionJSON string
	if err := chromedp.Run(bctx,
		chromedp.Navigate("https://artofproblemsolving.com/community/c3413"),
		chromedp.Sleep(8*time.Second),
		chromedp.Evaluate(`JSON.stringify(window.AoPS && window.AoPS.session ? window.AoPS.session : null)`, &sessionJSON),
	); err != nil {
		cancelBrowser()
		cancelAlloc()
		logger.Fatal("chromedp init: ", err)
	}
	if sessionJSON == "" || sessionJSON == "null" {
		cancelBrowser()
		cancelAlloc()
		logger.Fatal("AoPS.session not found after navigation; the Cloudflare challenge probably did not resolve")
	}
	data := ForumSession{
		browserCtx:    bctx,
		cancelBrowser: cancelBrowser,
		cancelAlloc:   cancelAlloc,
	}
	if err := json.Unmarshal([]byte(sessionJSON), &data); err != nil {
		cancelBrowser()
		cancelAlloc()
		logger.Fatal("parsing AoPS.session: ", err)
	}
	logger.Printf("Forum session ready (user_id=%d, logged_in=%v)", data.UserId, data.LoggedIn)
	return data
}

func (f *ForumSession) Close() {
	if f.cancelBrowser != nil {
		f.cancelBrowser()
	}
	if f.cancelAlloc != nil {
		f.cancelAlloc()
	}
}

// resolveChallenge re-navigates the browser to a community page so Cloudflare
// can issue a fresh cf_clearance cookie. Takes the write lock — all in-flight
// fetches finish first, then no new fetches start until navigation completes.
// Dedupes by wall-clock: if another goroutine just resolved, skip.
func (f *ForumSession) resolveChallenge() {
	resolveMu.Lock()
	defer resolveMu.Unlock()
	// Dedupe window matches the empirical wall duration (~22s). If a previous
	// goroutine already paused-and-re-navigated within the last 30s, the cookie
	// is fresh enough; skip.
	lastResolveAtMu.Lock()
	if !lastResolveAt.IsZero() && time.Since(lastResolveAt) < 30*time.Second {
		lastResolveAtMu.Unlock()
		return
	}
	lastResolveAtMu.Unlock()

	logger.Println("Cloudflare wall hit. Cooling down 30s, then re-navigating...")
	time.Sleep(30 * time.Second)
	ctx, cancel := context.WithTimeout(f.browserCtx, 45*time.Second)
	defer cancel()
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://artofproblemsolving.com/community/c3413"),
		chromedp.Sleep(6*time.Second),
	); err != nil {
		logger.Println("resolveChallenge:", err)
	}
	lastResolveAtMu.Lock()
	lastResolveAt = time.Now()
	lastResolveAtMu.Unlock()
}

const fetchAjaxJS = `(async () => {
	const r = await fetch('/m/community/ajax.php', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
			'X-Requested-With': 'XMLHttpRequest',
			'Accept': 'application/json, text/javascript, */*; q=0.01'
		},
		body: %s
	});
	return await r.text();
})()`

// postAjax posts to /m/community/ajax.php from inside the headless browser so
// the request inherits cf_clearance and the matching TLS handshake. Retries up
// to 3x with backoff on JS exceptions or non-JSON responses.
func (f *ForumSession) postAjax(body url.Values) ([]byte, error) {
	if body.Get("aops_logged_in") == "" {
		body.Set("aops_logged_in", strconv.FormatBool(f.LoggedIn))
		body.Set("aops_user_id", strconv.Itoa(f.UserId))
		body.Set("aops_session_id", f.SessionId)
	}
	bodyLit, err := json.Marshal(body.Encode())
	if err != nil {
		return nil, err
	}
	js := fmt.Sprintf(fetchAjaxJS, string(bodyLit))

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt*2) * time.Second)
		}
		// Rate cap first — paces total req/sec across all goroutines.
		reserveRequestSlot()
		// RLock: many fetches can run concurrently, but resolveChallenge can
		// take the write lock to pause everyone for re-navigation.
		resolveMu.RLock()
		requestSem <- struct{}{}
		ctx, cancel := context.WithTimeout(f.browserCtx, 5*time.Minute)
		var result string
		runErr := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			res, exc, e := runtime.Evaluate(js).WithAwaitPromise(true).WithReturnByValue(true).Do(ctx)
			if e != nil {
				return e
			}
			if exc != nil {
				return fmt.Errorf("js exception: %s", exc.Text)
			}
			return json.Unmarshal(res.Value, &result)
		}))
		cancel()
		<-requestSem
		resolveMu.RUnlock()
		if runErr != nil {
			lastErr = runErr
			continue
		}
		if json.Valid([]byte(result)) {
			return []byte(result), nil
		}
		// If Cloudflare bounced us with a fresh challenge, re-solve. The call
		// dedupes — if another goroutine resolved <10s ago, it's a no-op.
		if strings.Contains(result, "Just a moment...") {
			f.resolveChallenge()
		}
		preview := result
		if len(preview) > 160 {
			preview = preview[:160]
		}
		lastErr = fmt.Errorf("non-JSON response from AoPS via in-page fetch (preview: %s)", preview)
	}
	return nil, lastErr
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
	respbody, err := f.postAjax(url.Values{"a": {"fetch_topic"}, "topic_id": {strconv.Itoa(id)}})
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	serialized := TopicResponse{}
	err = json.Unmarshal(respbody, &serialized)
	if err != nil || serialized.Response.Topic == nil {
		serializederror := ErrorResponse{}
		if sererr := json.Unmarshal(respbody, &serializederror); sererr != nil {
			return nil, sererr
		}
		if len(serializederror.Code) > 0 {
			return nil, errors.New(serializederror.Code)
		}
		return nil, err
	}
	logger.Println("Finished Parsing Forum Topic", id)
	return &serialized, nil
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
	respbody, err := f.postAjax(url.Values{"a": {"fetch_category_data"}, "category_id": {strconv.Itoa(id)}})
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	serialized := CategoryResponse{}
	err = json.Unmarshal(respbody, &serialized)
	if err != nil || serialized.Response.Category == nil {
		serializederror := ErrorResponse{}
		if sererr := json.Unmarshal(respbody, &serializederror); sererr != nil {
			return nil, sererr
		}
		if len(serializederror.Code) > 0 {
			return nil, errors.New(serializederror.Code)
		}
		return nil, err
	}
	for i, x := range serialized.Response.Category.Items {
		// in the second case described above, this effectively does
		// nothing.
		r, perr := parseProblemRenderedHTML(x.PostData.Rendered)
		if perr != nil {
			logger.Println("parseProblemRenderedHTML:", perr)
			continue
		}
		serialized.Response.Category.Items[i].PostData.Canonical = r
	}
	logger.Println("Finished Parsing Forum Category", id)
	return &serialized, nil
}
