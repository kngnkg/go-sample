package store

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

// 疎通を確認する
func (r *Repository) Ping(ctx context.Context, db DBConnection) error {
	xdb, ok := db.(*sqlx.DB)
	if !ok {
		return errors.New("invalid args")
	}
	return xdb.PingContext(ctx)
}
