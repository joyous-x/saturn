package model

import (
	"context"
	"database/sql"
)

// Transaction ...
func Transaction(ctx context.Context, db *sql.DB, opt *sql.TxOptions, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, opt)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = fn(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}
