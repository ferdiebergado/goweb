//go:generate mockgen -destination=mock/dbrepo_mock.go -package=mock . DBRepo
package repository

import (
	"context"
	"database/sql"
)

type DBRepo interface {
	DBVersion(ctx context.Context) (string, error)
}

type dbRepo struct {
	db *sql.DB
}

func NewDBRepo(db *sql.DB) DBRepo {
	return &dbRepo{
		db: db,
	}
}

const VersionQuery = "SELECT version()"

func (r *dbRepo) DBVersion(ctx context.Context) (string, error) {
	var v string
	if err := r.db.QueryRowContext(ctx, VersionQuery).Scan(&v); err != nil {
		return "", err
	}

	return v, nil
}
