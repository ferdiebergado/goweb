//go:generate mockgen -destination=mock/repository_mock.go -package=mock . Repository
package repository

import (
	"context"
	"database/sql"
)

type Repository interface {
	PingContext(ctx context.Context) error
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) PingContext(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
