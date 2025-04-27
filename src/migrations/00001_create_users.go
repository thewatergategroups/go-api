package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)


func init() {
	goose.AddMigrationContext(upCreateUsers, downCreateUsers)
}

func upCreateUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT now()
		);
	`)
	return err
}

func downCreateUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE users;`)
	return err
}