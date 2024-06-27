package db

import (
	"context"
	"fmt"

	"github.com/karust/openserp/core"
	"github.com/karust/openserp/ent"
)

func InsertBulk(ctx context.Context, db *ent.Client, results []core.SearchResult, loc, lang, searchQ string, sqID int) error {
	tx, err := db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("starting a transaction: %w", err)
	}

	b := make([]*ent.SERPCreate, 0, len(results))
	for _, result := range results {
		b = append(b,
			tx.SERP.Create().
				SetTitle(result.Title).
				SetDescription(result.Description).
				SetURL(result.URL).
				SetKeyWords(nil2Zero(result.KeyWords)).
				SetEmails(nil2Zero(result.Emails)).
				SetPhones(nil2Zero(result.Phones)).
				SetSqID(sqID),
		)
	}
	_, err = db.SERP.CreateBulk(b...).Save(ctx)
	if err != nil {
		return rollback(tx, err)
	}
	return tx.Commit()
}

func nil2Zero(s []string) []string {
	if len(s) == 0 {
		return make([]string, 0)
	}
	return s
}
