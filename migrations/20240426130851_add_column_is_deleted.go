package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddColumnIsDeleted, downAddColumnIsDeleted)
}

func upAddColumnIsDeleted(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "ALTER TABLE urls ADD COLUMN is_deleted bool not null default false")
	return err
}

func downAddColumnIsDeleted(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "ALTER TABLE urls DROP COLUMN is_deleted")
	return err
}
