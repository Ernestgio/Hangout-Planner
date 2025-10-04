package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/domain"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/repository"
)

func TestNewUserRepository(t *testing.T) {
	db, _ := setupDB(t)
	repo := repository.NewUserRepository(db)
	require.NotNil(t, repo)
}

func TestCreateUser_Success(t *testing.T) {
	db, mock := setupDB(t)
	repo := repository.NewUserRepository(db)

	user := &domain.User{
		Name:     "Ernest",
		Email:    "ernest@example.com",
		Password: "hashed_password",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users` (`id`,`name`,`email`,`password`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?)").
		WithArgs(sqlmock.AnyArg(), user.Name, user.Email, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateUser(user)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_Error(t *testing.T) {
	db, mock := setupDB(t)
	repo := repository.NewUserRepository(db)
	dbError := errors.New("db error")

	user := &domain.User{
		Name:     "Ernest",
		Email:    "ernest@example.com",
		Password: "hashed_password",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users` (`id`,`name`,`email`,`password`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?)").
		WithArgs(sqlmock.AnyArg(), user.Name, user.Email, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnError(dbError)
	mock.ExpectRollback()

	err := repo.CreateUser(user)
	require.Error(t, err)
	require.Equal(t, dbError, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_Success(t *testing.T) {
	db, mock := setupDB(t)
	repo := repository.NewUserRepository(db)

	id := uuid.New()
	now := time.Now()
	email := "found@example.com"

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at", "deleted_at"}).
		AddRow(id, "Found User", email, "pwd", now, now, nil)

	expectedSQL := "SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?"
	mock.ExpectQuery(expectedSQL).WithArgs(email, 1).WillReturnRows(rows)

	user, err := repo.GetUserByEmail(email)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, email, user.Email)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, mock := setupDB(t)
	repo := repository.NewUserRepository(db)
	email := "notfound@example.com"

	expectedSQL := "SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?"
	mock.ExpectQuery(expectedSQL).WithArgs(email, 1).WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.GetUserByEmail(email)
	require.Error(t, err)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	require.Nil(t, user)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_DBError(t *testing.T) {
	db, mock := setupDB(t)
	repo := repository.NewUserRepository(db)
	email := "error@example.com"
	dbError := errors.New("db error")

	expectedSQL := "SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?"
	mock.ExpectQuery(expectedSQL).WithArgs(email, 1).WillReturnError(dbError)

	user, err := repo.GetUserByEmail(email)
	require.Error(t, err)
	require.Equal(t, dbError, err)
	require.Nil(t, user)
	require.NoError(t, mock.ExpectationsWereMet())
}
