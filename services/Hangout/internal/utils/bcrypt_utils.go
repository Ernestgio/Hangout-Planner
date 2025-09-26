package utils

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type BcryptUtils interface {
	CompareHashAndPassword(hashedPassword, password string) error
	GenerateFromPassword(password string) (string, error)
}

type bcryptUtils struct {
	cost int
}

func NewBcryptUtils(cost int) BcryptUtils {
	return &bcryptUtils{
		cost: cost,
	}
}

func (b *bcryptUtils) CompareHashAndPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return apperrors.ErrInvalidCredentials
		}
		return err
	}
	return nil
}

func (b *bcryptUtils) GenerateFromPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
