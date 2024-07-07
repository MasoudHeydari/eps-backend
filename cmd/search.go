package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/karust/openserp/db"

	"github.com/karust/openserp/baidu"
	"github.com/karust/openserp/core"
	"github.com/karust/openserp/google"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var searchCMD = &cobra.Command{
	Use:     "search",
	Aliases: []string{"find"},
	Short:   "Search results using chosen web search engine (google, yandex, baidu)",
	Run:     search,
}

func search(_ *cobra.Command, _ []string) {
	fmt.Println("time is: ", time.Now().Round(time.Second))
}

func search02(_ *cobra.Command, _ []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*2)
	defer cancel()
	client, err := db.NewDB()
	if err != nil {
		logrus.Fatalf("failed to connect to DB, error: %v", err)
	}
	searchQueries, err := db.GetAllSearchQueries(ctx, client)
	if err != nil {
		logrus.Error("failed to search: %w", err)
		os.Exit(1)
	}

	for _, searchQuery := range searchQueries {
		defer func() {
			if recErr := recover(); recErr != nil {
				logrus.Infof("search.Recover.err: %v", recErr)
			}
		}()
		go func(sq db.SearchQuery) {
			search01(ctx, client, sq.Location, sq.Language, sq.Query, sq.Id)
		}(searchQuery)
	}
	<-ctx.Done()
	logrus.Info("crawler exit successfully")
}

func searchBrowser(engine core.SearchEngine, query core.Query) ([]core.SearchResult, error) {
	return engine.Search(query)
}

//func searchRaw(engineType string, query core.Query) ([]core.SearchResult, error) {
//	logrus.Warn("Browserless results are very inconsistent or may not even work!")
//
//	switch strings.ToLower(engineType) {
//	case "yandex":
//		return yandex.Search(query)
//	case "google":
//		return google.Search(query)
//	case "baidu":
//		return baidu.Search(query)
//	default:
//		logrus.Infof("No `%s` search engine found", engineType)
//	}
//	return nil, nil
//}

func buildEngine(engineType string) core.SearchEngine {
	opts := core.BrowserOpts{
		IsHeadless:    !config.App.IsBrowserHead, // Disable headless if browser head mode is set
		IsLeakless:    config.App.IsLeakless,
		Timeout:       time.Second * time.Duration(config.App.Timeout),
		LeavePageOpen: config.App.IsLeaveHead,
	}

	if config.App.IsDebug {
		opts.IsHeadless = false
	}

	browser, err := core.NewBrowser(opts)
	if err != nil {
		logrus.Error(err)
	}
	// browser.Initialize()
	var engine core.SearchEngine
	switch strings.ToLower(engineType) {
	case "google":
		engine = google.New(*browser, config.GoogleConfig)
	case "baidu":
		engine = baidu.New(*browser, config.BaiduConfig)
	default:
	}
	return engine
}

func init() {
	RootCmd.AddCommand(searchCMD)
	searchCMD.Flags().IntP("offset", "o", 0, "set offset for your search query")
}
