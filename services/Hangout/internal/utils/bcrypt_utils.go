package utils

import (
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"golang.org/x/crypto/bcrypt"
)

func CompareHashAndPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return apperrors.ErrInvalidCredentials
		}
		return err
	}
	return nil
}
