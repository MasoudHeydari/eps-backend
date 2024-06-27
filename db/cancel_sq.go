package db

import (
	"context"

	"github.com/karust/openserp/ent"
)

func CancelSQ(ctx context.Context, db *ent.Client, sqID int) error {
	_, err := db.SearchQuery.UpdateOneID(sqID).
		SetIsCanceled(true).
		Save(ctx)
	return err
}
