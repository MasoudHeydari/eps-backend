package db

import (
	"context"

	"github.com/karust/openserp/core"
	"github.com/karust/openserp/ent"
	"github.com/sirupsen/logrus"
)

func InsertBulk(ctx context.Context, db *ent.Client, results []core.SearchResult, loc, lang, searchQ string, sqID int) error {
	logrus.Println("InsertBulk: len of results to insert: ", len(results))
	lenOfInsertedRows := 0
	for _, result := range results {
		_, err := db.SERP.Create().
			SetTitle(result.Title).
			SetDescription(result.Description).
			SetURL(result.URL).
			SetKeyWords(nil2Zero(result.KeyWords)).
			SetEmails(nil2Zero(result.Emails)).
			SetPhones(nil2Zero(result.Phones)).
			SetSqID(sqID).
			Save(ctx)
		if err != nil {
			logrus.Info("InsertBulk.SERP.Create: ", err)
			continue
		}
		lenOfInsertedRows++
	}
	logrus.Printf("InsertBulk: %d new rows added successfully\n", lenOfInsertedRows)
	return nil
}

func nil2Zero(s []string) []string {
	if len(s) == 0 {
		return make([]string, 0)
	}
	return s
}
