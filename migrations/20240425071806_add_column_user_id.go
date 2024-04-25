package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddColumnUserId, downAddColumnUserId)
}

func upAddColumnUserId(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "ALTER TABLE urls ADD COLUMN user_id uuid not null")
	return err
}

func downAddColumnUserId(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "ALTER TABLE urls DROP COLUMN user_id")
	return err
}
