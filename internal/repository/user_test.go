package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/radiatus-ai/auth-service/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock
}

func TestCreateUser(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "password",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"users\"").
		WithArgs(user.ID, user.Email, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmail(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)

	email := "test@example.com"
	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
		AddRow(userID, email, "password", "2023-01-01 00:00:00", "2023-01-01 00:00:00")

	mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email = \\$1").
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.GetByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExistsByEmail(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewUserRepository(db)

	email := "test@example.com"

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"users\" WHERE email = \\$1").
		WithArgs(email).
		WillReturnRows(rows)

	exists, err := repo.ExistsByEmail(email)
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}
