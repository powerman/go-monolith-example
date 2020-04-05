package migrations

import (
	"database/sql"
)

func init() {
	goose.AddMigration(upNoop, downNoop)
}

func upNoop(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func downNoop(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil // migrate.ErrDownNotSupported
}
