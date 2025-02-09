package repository

import (
	"context"
	"database/sql"

	"github.com/ferdiebergado/goweb/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

const CreateUserQuery = `
INSERT INTO users (email, password_hash)
VALUES $1, $2
RETURNING id, email, created_at, updated_at
`

func (r *UserRepo) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	var newUser model.User
	if err := r.db.QueryRowContext(ctx, CreateUserQuery, user.Email, user.PasswordHash).
		Scan(&newUser.ID, &newUser.Email, &newUser.CreatedAt, &newUser.UpdatedAt); err != nil {
		return nil, err
	}
	return &newUser, nil
}
