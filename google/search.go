package google

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func (gogl *Google) findTotalResults(page *rod.Page) (int, error) {
	resultsStats, err := page.Timeout(gogl.GetSelectorTimeout()).Search("div#result-stats")
	if err != nil {
		return 0, errors.New("Result stats not found: " + err.Error())
	}

	stats, err := resultsStats.First.Text()
	if err != nil {
		return 0, errors.New("Cannot extract result stats text: " + err.Error())
	}

	// Escape moment with `seconds` and extract digits
	allNums := gogl.findNumRgxp.FindAllString(stats[:len(stats)-15], -1)
	stats = strings.Join(allNums, "")

	total, err := strconv.Atoi(stats)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (gogl *Google) isCaptcha(page *rod.Page) bool {
	_, err := page.Timeout(gogl.GetSelectorTimeout()).Search("form#captcha-form")
	if err != nil {
		return false
	}
	return true
}

func (gogl *Google) preparePage(page *rod.Page) {
	// Remove "similar queries" lists
	page.Eval(";(() => { document.querySelectorAll(`div[data-initq]`).forEach( el => el.remove());  })();")
}

func (gogl *Google) Search(query core.Query) ([]core.SearchResult, error) {
	logrus.Tracef("Start Google search, query: %+v", query)

	var searchResults []core.SearchResult

	// Build URL from query struct to open in browser
	url, err := BuildURL(query)
	if err != nil {
		return nil, err
	}

	page := gogl.Navigate(url)
	gogl.preparePage(page)

	results, err := page.Timeout(gogl.Timeout).Search("div[data-hveid][data-ved][lang], div[data-surl][jsaction]")
	if err != nil {
		defer page.Close()
		logrus.Errorf("Cannot parse search results: %s", err)
		return nil, core.ErrSearchTimeout
	}

	// Check why no results, maybe captcha?
	if results == nil {
		defer page.Close()

		if gogl.isCaptcha(page) {
			logrus.Errorf("Google captcha occurred during: %s", url)
			return nil, core.ErrCaptcha
		}
		return nil, err
	}

	totalResults, err := gogl.findTotalResults(page)
	if err != nil {
		logrus.Errorf("Error capturing total results: %v", err)
	}
	logrus.Infof("%d total results found", totalResults)

	resultElements, err := results.All()
	if err != nil {
		return nil, err
	}

	for i, r := range resultElements {
		// Get URL
		link, err := r.Element("a")
		if err != nil {
			continue
		}
		linkText, err := link.Property("href")
		if err != nil {
			logrus.Error("No `href` tag found")
		}

		// Get title
		titleTag, err := link.Element("h3")
		if err != nil {
			logrus.Error("No `h3` tag found")
			continue
		}

		title, err := titleTag.Text()
		if err != nil {
			logrus.Error("Cannot extract text from title")
			title = "No title"
		}

		// Get description
		// doesn't catch all
		descTag, err := r.Element(`div[data-sncf~="1"]`)
		desc := ""
		if err != nil {
			logrus.Trace(`No description 'div[data-sncf~="1"]' tag found`)
		} else {
			desc = descTag.MustText()
		}

		// extract contact-info
		emails, err := gogl.extractEmails(linkText.String())
		if err != nil {
			logrus.Errorf("Search: %v", err)
		}

		phones, err := gogl.extractPhoneNumbersFromAllPossibleURLs(linkText.String())
		if err != nil {
			logrus.Errorf("Search: %v", err)
		}

		// extract key-words
		var keyWords []string
		keyWords, err = gogl.extractKeywords(linkText.String())
		if err != nil {
			logrus.Errorf("Search: %v", err)
		}

		gR := core.SearchResult{
			Rank:        i + 1,
			URL:         linkText.String(),
			Title:       title,
			Phones:      phones,
			Emails:      emails,
			KeyWords:    keyWords,
			Description: desc,
		}
		searchResults = append(searchResults, gR)
	}

	if !gogl.Browser.LeavePageOpen {
		err = page.Close()
		if err != nil {
			logrus.Error(err)
		}
	}

	return searchResults, nil
}

func (gogl *Google) extractKeywords(path string) ([]string, error) {
	var (
		maxLen               = 6
		maxLenForEachKeyword = 72
	)
	keyWordsMap := make(map[string]struct{})
	keyWords := make([]string, 0)
	gogl.keywordCollector.OnHTML("h1, h2, h3", func(e *colly.HTMLElement) {
		keyWordsMap[e.Text] = struct{}{}
	})
	err := gogl.keywordCollector.Visit(path)
	if err != nil {
		return keyWords, fmt.Errorf("extractKeywords: %w", err)
	}
	if len(keyWordsMap) == 0 {
		return keyWords, fmt.Errorf("no keyword found")
	}
	for k := range keyWordsMap {
		k = strings.TrimSpace(k)
		k = strings.Trim(k, "\n")
		k = strings.Trim(k, "\t")
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
	for k := range emailsMap {
		emails = append(emails, k)
	}
	return emails, nil
}

func (gogl *Google) extractPhoneNumbers(path string) []string {
	matches := make([]string, 0)
	gogl.phoneCollector.OnHTML("div", func(e *colly.HTMLElement) {
		matches = append(
			matches,
			gogl.findPhoneRgxp.FindAllString(e.Text, -1)...,
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
		return phoneNums, fmt.Errorf("no phone number found")
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
