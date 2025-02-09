package repository_test

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/ferdiebergado/goweb/internal/model"
	"github.com/ferdiebergado/goweb/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := model.User{
		Email:        "abc@example.com",
		PasswordHash: "hashed",
	}

	mock.ExpectQuery(repository.CreateUserQuery).
		WithArgs(user.Email, user.PasswordHash).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "created_at", "updated_at"}).
			AddRow("1", user.Email, time.Now(), time.Now()))

	repo := repository.NewUserRepository(db)
	newUser, err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.NotNil(t, newUser, "New user should not be empty")
	assert.Equal(t, user.Email, newUser.Email, "emails should match")
	assert.NotZero(t, newUser.CreatedAt)
	assert.NotZero(t, newUser.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}
