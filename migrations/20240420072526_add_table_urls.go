package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddTableUrls, downAddTableUrls)
}

func upAddTableUrls(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "CREATE TABLE urls (key text PRIMARY KEY, url text UNIQUE)")
	return err
}

func downAddTableUrls(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE urls")
	return err
}
