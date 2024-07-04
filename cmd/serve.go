package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/karust/openserp/core"
	"github.com/karust/openserp/db"
	"github.com/karust/openserp/ent"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const googleEngine = "google"

type Server struct {
	db *ent.Client
}

var serveCMD = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"listen"},
	Short:   "Start HTTP server, to provide search engine results via API",
	Args:    cobra.MatchAll(cobra.NoArgs),
	Run:     serve,
}

func serve(cmd *cobra.Command, args []string) {
	//opts := core.BrowserOpts{
	//	IsHeadless:    !config.App.IsBrowserHead, // Disable headless if browser head mode is set
	//	IsLeakless:    config.App.IsLeakless,
	//	Timeout:       time.Second * time.Duration(config.App.Timeout),
	//	LeavePageOpen: config.App.IsLeaveHead,
	//}
	//
	//if config.App.IsDebug {
	//	opts.IsHeadless = false
	//}
	//
	//browser, err := core.NewBrowser(opts)
	//if err != nil {
	//	logrus.Error(err)
	//}
	//
	//yand := yandex.New(*browser, config.YandexConfig)
	//gogl := google.New(*browser, config.GoogleConfig)
	//baidu := baidu.New(*browser, config.BaiduConfig)
	//
	//serv := core.NewServer(config.App.Host, config.App.Port, gogl, yand, baidu)
	//serv.Listen()
	//client, err := db.NewDB()
	//if err != nil {
	//	logrus.Errorf("Failed to connect to DB, error: %v", err)
	//	return
	//}

	logrus.Fatal(Start())
}

func init() {
	RootCmd.AddCommand(serveCMD)
}

type searchQ struct {
	Language string `json:"lang"`
	Location string `json:"loc"`
	Query    string `json:"q"`
}

type CancelSQ struct {
	SQID int `json:"sq_id"` // SQID is Search Query ID
}

type GetAllSearchResults struct {
	SQID int `param:"sq_id"` // SQID is Search Query ID
	Page int `query:"page"`
}

type ExportCSV struct {
	SQID int `param:"sq_id"` // SQID is Search Query ID
}

func Start() error {
	client, err := db.NewDB()
	if err != nil {
		return fmt.Errorf("failed to connect to DB, error: %v", err)
	}
	server := Server{db: client}
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.POST("/api/v1/search", server.Search)
	e.GET("/api/v1/search/:sq_id", server.GetSearchResults)
	e.PATCH("/api/v1/search", server.CancelSearchQuery)
	e.GET("/api/v1/export/:sq_id", server.ExportCSV)
	e.GET("/api/v1/search", server.GetAllSearchQueries)
	// go server.searchInBackground()
	return e.Start(":9999")
}

func (s *Server) Search(c echo.Context) error {
	//  curl -X POST -w "%{http_code}\n" http://localhost:9999/api/v1/search -H "Content-Type: application/json" -d '{"loc": "NL", "lang": "En", "q": "Golang"}'
	sq := new(searchQ)
	if err := c.Bind(sq); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	sqID, err := db.InsertNewSeaerchQuery(c.Request().Context(), s.db, sq.Location, sq.Language, sq.Query)
	if err != nil {
		switch {
		case ent.IsConstraintError(err):
			return c.JSON(http.StatusConflict, err)
		default:
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("recover goroutine panic: %v", err)
			}
		}()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
		defer cancel()
		search01(ctx, s.db, sq.Location, sq.Language, sq.Query, sqID)
		logrus.Printf("Search Query: %+v finished\n", sq)
	}()

	return c.JSON(http.StatusOK, echo.Map{"sq_id": sqID})
}

func (s *Server) GetSearchResults(c echo.Context) error {
	dto := new(GetAllSearchResults)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if dto.SQID == 0 {
		return c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	serps, err := db.GetAllResult(c.Request().Context(), s.db, dto.SQID, dto.Page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, serps)
}

func (s *Server) ExportCSV(c echo.Context) error {
	dto := new(ExportCSV)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if dto.SQID == 0 {
		return c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	csvAbsPatch, fileName, err := db.ExportCSV(c.Request().Context(), s.db, dto.SQID)
	if err != nil {
		logrus.Info("ExportCSV: ", err)
		return c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	return c.Attachment(csvAbsPatch, fileName)
}

func (s *Server) CancelSearchQuery(c echo.Context) error {
	// curl -X PATCH -w "%{http_code}\n" http://localhost:9999/api/v1/search -H "Content-Type: application/json" -d '{"sq_id" : 1}'
	dto := new(CancelSQ)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if dto.SQID == 0 {
		return c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	err := db.CancelSQ(c.Request().Context(), s.db, dto.SQID)
	if err != nil {
		logrus.Info("CancelSearchQuery.CancelSQ: %w", err)
		switch {
		case ent.IsNotFound(err):
			return c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		default:
			return c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) GetAllSearchQueries(c echo.Context) error {
	sqs, err := db.GetAllSearchQueries(c.Request().Context(), s.db)
	if err != nil {
		logrus.Info("GetAllSearchQueries: %w", err)
		c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, echo.Map{"search_queries": sqs})
}

func search01(ctx context.Context, client *ent.Client, loc, lang, searchQ string, sqID int) {
	query := core.Query{
		Text:     searchQ,
		LangCode: lang,
		Location: loc,
		Limit:    10,
		Offset:   0,
	}
	engine := buildEngine(googleEngine)
	if engine == nil {
		logrus.Errorf("Failed to build Engine, No `%s` search engine found", googleEngine)
		return
	}

	for {
		select {
		case <-ctx.Done():
			logrus.Println("Inserted all found results")
			return
		default:
			results, err := searchBrowser(engine, query)
			if err != nil {
				logrus.Error(err)
				//return
				continue
			}

			// Save found records into the DB
			err = db.InsertBulk(ctx, client, results, loc, lang, searchQ, sqID)
			if err != nil {
				switch {
				case ent.IsConstraintError(err):
				default:
					logrus.Errorf("failed to insert results to DB, error: %v", err)
					return
				}
			}
			query.NextPage()
		}
	}
}

func (s *Server) searchInBackground() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Hour*2)
		defer cancel()
		searchQueries, err := db.GetAllSearchQueries(ctx, s.db)
		if err != nil {
			logrus.Infof("searchInBackground: %v\n", err)
			return
		}
		for _, searchQuery := range searchQueries {
			time.Sleep(10 * time.Second)
			go func(sq db.SearchQuery) {
				defer func() {
					if recErr := recover(); recErr != nil {
						logrus.Infof("searchInBackground.Recover.err: %v\n", recErr)
					}
				}()
				logrus.Infof("crawler just resumed for Search Query: %v\n", sq)
				search01(ctx, s.db, sq.Location, sq.Language, sq.Query, sq.Id)
			}(searchQuery)
		}
		<-ctx.Done()
		logrus.Info("crawler sleeped for one day")
		time.Sleep(24 * time.Hour)
	}
}
