package db

import (
	"context"
	"fmt"

	"github.com/karust/openserp/ent"
	"github.com/karust/openserp/ent/serp"
)

func GetAllResult(ctx context.Context, db *ent.Client, sqID, page int) ([]SERP, error) {
	offSet := page * 5
	entSERPs, err := db.SERP.Query().
		Where(
			serp.SqID(sqID),
			serp.IsRead(false),
		).
		Offset(offSet).
		Limit(5).
		Order(ent.Desc(serp.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]SERP, 0, len(entSERPs))
	for _, entSERP := range entSERPs {
		fmt.Println("kw: ", entSERP.KeyWords)
		results = append(results,
			SERP{
				URL:         entSERP.URL,
				Title:       entSERP.Title,
				Description: entSERP.Description,
				Phones:      entSERP.Phones,
				Emails:      entSERP.Emails,
				Keywords:    entSERP.KeyWords,
				IsRead:      entSERP.IsRead,
				CreatedAt:   entSERP.CreatedAt,
			},
		)
	}
	return results, nil
}
