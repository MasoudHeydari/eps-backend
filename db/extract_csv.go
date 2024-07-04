package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/karust/openserp/ent"
	"github.com/karust/openserp/ent/searchquery"
)

const (
	csvAbsPathInPostgresContainer = "/tmp/eps/db/data/export.csv"
	csvAbsPathInEPSContainer      = "/tmp/eps/db/csv/export.csv"
)

func ExportCSV(ctx context.Context, db *ent.Client, sqID int) (csvAbsFilePath, fileName string, err error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return "", "", fmt.Errorf("starting a transaction: %w", err)
	}
	entSq, err := tx.SearchQuery.Query().Where(searchquery.ID(sqID)).First(ctx)
	if err != nil {
		return "", "", rollback(tx, err)
	}
	// basePath, err := os.Getwd()
	// if err != nil {
	// 	return "", rollback(tx, fmt.Errorf("getwd: %w", err))
	// }
	// fmt.Println(entSq.Query)
	csvFileName := fmt.Sprintf("%s-%s.csv", entSq.Query, time.Now().Format("2006-01-02_15:04:05"))
	// csvAbsFilePath = "/tmp/eps/db/data/export.csv" // "/tmp/eps/db/csv/export.csv"
	// csvAbsFilePathInPsqlContainer := fmt.Sprintf("/tmp/eps/%s", csvFileName)
	// csvFile, err := os.Create(csvAbsFilePath)
	// if err != nil {
	// 	return "", rollback(tx, fmt.Errorf("openFile: %w", err))
	// }
	// defer func() {
	// 	csvFile.Close()
	// 	if err != nil {
	// 		// _ = os.Remove(csvFile.Name())
	// 	}
	// }()
	q := fmt.Sprintf(`COPY (
		SELECT
			serps.url,
			serps.title,
			serps.description,
			serps.phones,
			serps.emails,
			serps.key_words,
			search_queries.location,
			search_queries.language,
			serps.created_at
		FROM serps
		JOIN search_queries
		ON serps.sq_id=search_queries.id
		WHERE serps.sq_id=%d
		ORDER BY serps.id)
	TO '%s' WITH (FORMAT CSV, HEADER);`, sqID, csvAbsPathInPostgresContainer)

	_, err = tx.ExecContext(ctx, q) //, csvAbsFilePath)
	if err != nil {
		return "", "", rollback(tx, fmt.Errorf("ExecContext: %w", err))
	}
	err = tx.Commit()
	if err != nil {
		return "", "", err
	}
	return csvAbsPathInEPSContainer, csvFileName, nil
}

func rollback(tx *ent.Tx, err error) error {
	fmt.Println("rollback: error: ", err)
	if rErr := tx.Rollback(); rErr != nil {
		err = fmt.Errorf("%w: %v", err, rErr)
	}
	return err
}
