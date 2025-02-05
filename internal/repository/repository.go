//go:generate mockgen -destination=mock/repository_mock.go -package=mock . Repository
package repository

import (
	"context"
	"database/sql"
)

type Repository interface {
	Ping(ctx context.Context) error
}

type repo struct {
	db *sql.DB
}

var _ Repository = (*repo)(nil)

func NewRepository(db *sql.DB) Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
