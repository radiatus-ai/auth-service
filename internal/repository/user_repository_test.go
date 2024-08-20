package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/radiatus-ai/auth-service/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

func TestCreateUser_Success(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repo := NewUserRepository(db)
	user := &model.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		GoogleID: "google123",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users" (.+) VALUES (.+)`).
		WithArgs(user.ID, user.Email, user.GoogleID, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(user)

	assert.NoError(t, err)
	// assert.Equal(t, user, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
