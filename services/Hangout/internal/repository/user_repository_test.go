package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/models"
	"github.com/Ernestgio/Hangout-Planner/services/Hangout/internal/repository"
)

func setupDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)
	return db, mock
}

func TestCreateUser(t *testing.T) {
	db, mock := setupDB(t)
	repo := repository.NewUserRepository(db)

	u := &models.User{
		ID:       uuid.New(),
		Name:     "X",
		Email:    "x@example.com",
		Password: "secret",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users` (`id`,`name`,`email`,`password`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?)").
		WithArgs(
			u.ID, u.Name, u.Email, u.Password,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	require.NoError(t, repo.CreateUser(u))
}

func TestGetUserByEmail(t *testing.T) {
	db, mock := setupDB(t)
	repo := repository.NewUserRepository(db)

	id := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "created_at", "updated_at", "deleted_at",
	}).AddRow(id, "Y", "y@example.com", "pwd", now, now, nil)

	mock.ExpectQuery("SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
		WithArgs("y@example.com", 1).
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail("y@example.com")
	require.NoError(t, err)
	require.Equal(t, "y@example.com", user.Email)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, mock := setupDB(t)
	repo := repository.NewUserRepository(db)

	mock.ExpectQuery("SELECT .* FROM users.*").
		WithArgs("missing@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.GetUserByEmail("missing@example.com")
	require.Error(t, err)
	require.Nil(t, user)

	require.Error(t, mock.ExpectationsWereMet())
}
