//go:generate mockgen -destination=mock/user_repo_mock.go -package=mock . UserRepo
package repository

import (
	"context"
	"database/sql"

	"github.com/ferdiebergado/goweb/internal/model"
)

type UserRepo interface {
	CreateUser(ctx context.Context, params CreateUserParams) (*model.User, error)
}

type userRepo struct {
	db *sql.DB
}

var _ UserRepo = (*userRepo)(nil)

func NewUserRepository(db *sql.DB) UserRepo {
	return &userRepo{db: db}
}

type CreateUserParams struct {
	Email        string
	PasswordHash string
}

const CreateUserQuery = `
INSERT INTO users (email, password_hash)
VALUES $1, $2
RETURNING id, email, created_at, updated_at
`

func (r *userRepo) CreateUser(ctx context.Context, params CreateUserParams) (*model.User, error) {
	var newUser model.User
	if err := r.db.QueryRowContext(ctx, CreateUserQuery, params.Email, params.PasswordHash).
		Scan(&newUser.ID, &newUser.Email, &newUser.CreatedAt, &newUser.UpdatedAt); err != nil {
		return nil, err
	}
	return &newUser, nil
}
