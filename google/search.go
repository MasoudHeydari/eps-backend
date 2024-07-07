package google

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	fakeUserAgent "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/gocolly/colly"
	"github.com/karust/openserp/core"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type Google struct {
	core.Browser
	core.SearchEngineOptions
	phoneCollector                             *colly.Collector
	keywordCollector                           *colly.Collector
	findNumRgxp, findPhoneRgxp, findEmailRegex *regexp.Regexp
	client                                     *http.Client
	limiter                                    *rate.Limiter
}

func New(browser core.Browser, opts core.SearchEngineOptions) *Google {
	phoneCollector := colly.NewCollector()
	phoneCollector.SetRequestTimeout(20 * time.Second)

	keywordCollector := colly.NewCollector()
	keywordCollector.SetRequestTimeout(20 * time.Second)
	gogl := Google{
		Browser:          browser,
		phoneCollector:   phoneCollector,
		keywordCollector: keywordCollector,
		client:           &http.Client{Timeout: 20 * time.Second},
		limiter:          rate.NewLimiter(1, 1),
	}
	opts.Init()
	gogl.SearchEngineOptions = opts

	gogl.findNumRgxp = regexp.MustCompile("\\d")
	gogl.findPhoneRgxp = regexp.MustCompile(`(?:[-+() ]*\d){10,13}`)
	gogl.findEmailRegex = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	return &gogl
}

func (gogl *Google) Name() string {
	return "google"
}

func (gogl *Google) GetRateLimiter() *rate.Limiter {
	ratelimit := rate.Every(gogl.GetRatelimit())
	return rate.NewLimiter(ratelimit, gogl.RateBurst)
}

func (gogl *Google) preparePage(page *rod.Page) {
	// Remove "similar queries" lists
	page.Eval(";(() => { document.querySelectorAll(`div[data-initq]`).forEach( el => el.remove());  })();")
}

func (gogl *Google) Search(query core.Query) ([]core.SearchResult, error) {
	if !gogl.limiter.Allow() {
		return nil, fmt.Errorf("too may request")
	}
	logrus.Tracef("Start Google search, query: %+v", query)

	var searchResults []core.SearchResult

	// Build URL from query struct to open in browser
	url, err := BuildURL(query)
	if err != nil {
		return nil, err
	}
	fmt.Println("Navigate to: ", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	req.Header.Set("User-Agent", fakeUserAgent.Random())
	res, err := gogl.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer time.Sleep(time.Second * time.Duration(getRandomSec()))
	defer gogl.client.CloseIdleConnections()
	defer res.Body.Close()
	fmt.Println("status code: ", res.StatusCode, " - Text: ", http.StatusText(res.StatusCode))
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	querySelector := "div[data-hveid][data-ved][lang], div[data-surl][jsaction]"
	doc.Find(querySelector).Each(func(i int, result *goquery.Selection) {
		link, found := result.Find("a").First().Attr("href")
		if !found {
			logrus.Info(
				"No url found",
				"query", query,
			)
			return
		}
		title := result.Find("h3").First().Text()
		desc := result.Find(`div[data-sncf~="1"]`).First().Text()
		emails, err := gogl.extractEmails(link)
		if err != nil {
			logrus.Errorf("Search: %v", err)
		}
		phones, err := gogl.extractPhoneNumbersFromAllPossibleURLs(link)
		if err != nil {
			logrus.Errorf("Search: %v", err)
		}
		keyWords, err := gogl.extractKeywords(link)
		if err != nil {
			logrus.Errorf("Search: %v", err)
		}
		gR := core.SearchResult{
			Rank:        i + 1,
			URL:         link,
			Title:       title,
			Phones:      phones,
			Emails:      emails,
			KeyWords:    keyWords,
			Description: desc,
		}
		searchResults = append(searchResults, gR)
	})
	return searchResults, nil
}

func (gogl *Google) extractKeywords(path string) ([]string, error) {
	var (
		maxLen               = 6
		maxLenForEachKeyword = 72
		h1h2h3QuerySelector  = "h1, h2, h3"
	)
	keyWordsMap := make(map[string]struct{})
	keyWords := make([]string, 0)
	defer func() {
		gogl.keywordCollector.OnHTMLDetach(h1h2h3QuerySelector)
	}()
	gogl.keywordCollector.OnHTML(h1h2h3QuerySelector, func(e *colly.HTMLElement) {
		keyWordsMap[e.Text] = struct{}{}
	})
	err := gogl.keywordCollector.Visit(path)
	if err != nil {
		return keyWords, fmt.Errorf("extractKeywords: %w", err)
	}
	if len(keyWordsMap) == 0 {
		return keyWords, fmt.Errorf("no keyword found from %q", path)
	}
	for k := range keyWordsMap {
		k = strings.TrimSpace(k)
		// k = strings.Trim(k, "\n")
		// k = strings.Trim(k, "\t")
		k = strings.ReplaceAll(k, "\t", "")
		k = strings.ReplaceAll(k, "\n", "")
		if len(k) > maxLenForEachKeyword {
			k = k[:maxLenForEachKeyword]
		}
		keyWords = append(keyWords, strings.TrimSpace(k))
	}
	if len(keyWords) > maxLen {
		keyWords = keyWords[:maxLen]
	}
	return keyWords, nil
}

func (gogl *Google) extractEmails(path string) ([]string, error) {
	emailsMap := make(map[string]struct{})
	emails := make([]string, 0)
	resp, err := gogl.client.Get(path)
	if err != nil {
		return emails, fmt.Errorf("extractEmails.Get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return emails, fmt.Errorf("extractEmails.ReadAll: %w", err)
	}
	emailsArray := gogl.findEmailRegex.FindAllString(string(body), -1)
	for _, email := range emailsArray {
		emailsMap[email] = struct{}{}
	}
	if len(emailsMap) == 0 {
		return nil, fmt.Errorf("no email found from %q", path)
	}
	for k := range emailsMap {
		emails = append(emails, k)
	}
	return emails, nil
}

func (gogl *Google) extractPhoneNumbers(path string) []string {
	matches := make([]string, 0)
	divQuerySelector := "div"
	defer func() {
		gogl.phoneCollector.OnHTMLDetach(divQuerySelector)
	}()
	gogl.phoneCollector.OnHTML(divQuerySelector, func(e *colly.HTMLElement) {
		newMatches := gogl.findPhoneRgxp.FindAllString(e.Text, -1)
		matches = append(
			matches,
			newMatches...,
		)
	})
	gogl.phoneCollector.Visit(path)
	return matches
}

func (gogl *Google) extractPhoneNumbersFromAllPossibleURLs(p string) ([]string, error) {
	phoneNums := make([]string, 0)
	pathes := make(map[string]struct{}, 0)
	u, err := url.Parse(p)
	if err != nil {
		return phoneNums, fmt.Errorf("extractPhoneNumbersFromAllPossibleURLs.Parse: %w", err)
	}
	pathes[p] = struct{}{}
	pathes[u.Scheme+"://"+u.Host+"/about-us"] = struct{}{}
	pathes[u.Scheme+"://"+u.Host+"/about"] = struct{}{}
	pathes[u.Scheme+"://"+u.Host+"/contact-us"] = struct{}{}
	pathes[u.Scheme+"://"+u.Host+"/contact"] = struct{}{}

	for pp := range pathes {
		phoneNums = append(phoneNums, gogl.extractPhoneNumbers(pp)...)
	}
	if len(phoneNums) == 0 {
		return phoneNums, fmt.Errorf("no phone number found from %q", p)
	}
	phonesMap := make(map[string]struct{})
	for _, phone := range phoneNums {
		phone = strings.TrimSpace(phone)
		phonesMap[phone] = struct{}{}
	}
	results := make([]string, 0, len(phonesMap))
	for k := range phonesMap {
		results = append(results, k)
	}
	return results, nil
}

func getRandomSec() int {
	sec := rand.Intn(30) + 45
	logrus.Println("sleep for: ", sec)
	return sec
}
