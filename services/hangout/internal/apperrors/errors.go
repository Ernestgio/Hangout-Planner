package apperrors

import "errors"

// application errors
var ErrAppPortRequired = errors.New("APP_PORT is required")
var ErrDbPasswordRequired = errors.New("DB_PASSWORD required in production")

// http errors
var ErrInvalidPayload = errors.New("invalid payload")

// business errors

// auth & user
var ErrUserAlreadyExists = errors.New("user with that email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")

// hangout
var ErrForbidden = errors.New("forbidden")
