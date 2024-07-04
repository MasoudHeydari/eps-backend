package db

import (
	"context"

	"entgo.io/ent/dialect"
	"github.com/karust/openserp/ent"
	"github.com/karust/openserp/ent/searchquery"
	_ "github.com/lib/pq"
)

func NewDB() (*ent.Client, error) {
	client, err := ent.Open(dialect.Postgres, "host=eps_db port=5432 user=eps_user dbname=eps_db password=eps_password sslmode=disable")
	if err != nil {
		return nil, err
	}
	err = client.Schema.Create(context.Background())
	if err != nil {
		return nil, err
	}
	return client, nil //
}

func InsertNewSeaerchQuery(ctx context.Context, db *ent.Client, loc, lang, searchQ string) (int, error) {
	sqEnt, err := db.SearchQuery.Create().
		SetLocation(loc).
		SetLanguage(lang).
		SetQuery(searchQ).
		Save(ctx)
	if err != nil {
		switch {
		case ent.IsConstraintError(err):
			err = db.SearchQuery.Update().
				Where(searchquery.IsCanceled(true)).
				SetIsCanceled(false).
				Exec(ctx)
			if err != nil {
				return -1, err
			}
		default:
			return -1, err
		}
	}
	return sqEnt.ID, nil
}

func GetAllSearchQueries(ctx context.Context, db *ent.Client) ([]SearchQuery, error) {
	entSearchQueries, err := db.SearchQuery.Query().
		Where(searchquery.IsCanceled(false)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	searchQueries := make([]SearchQuery, 0, len(entSearchQueries))
	for _, entSearchQuery := range entSearchQueries {

		searchQueries = append(searchQueries, SearchQuery{
			Id:         entSearchQuery.ID,
			Query:      entSearchQuery.Query,
			Language:   entSearchQuery.Language,
			Location:   entSearchQuery.Location,
			IsCanceled: entSearchQuery.IsCanceled,
			CreatedAt:  entSearchQuery.CreatedAt,
		})
	}
	return searchQueries, nil
}
