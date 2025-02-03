package repository_test

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/ferdiebergado/goweb/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestRepository_PingDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPing()

	repo := repository.NewRepository(db)
	err = repo.PingContext(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
