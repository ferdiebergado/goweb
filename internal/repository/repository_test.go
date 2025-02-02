package repository_test

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/ferdiebergado/goweb/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestDBVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(repository.VersionQuery).WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("1"))

	dbRepo := repository.NewDBRepo(db)
	v, err := dbRepo.DBVersion(context.Background())
	assert.NoError(t, err)
	assert.NotZero(t, v)
	assert.NoError(t, mock.ExpectationsWereMet())
}
