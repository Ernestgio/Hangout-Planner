package apperrors

import "errors"

// application errors
var ErrAppPortRequired = errors.New("APP_PORT is required")
var ErrDbPasswordRequired = errors.New("DB_PASSWORD required in production")

// http errors
var ErrInvalidPayload = errors.New("invalid payload")

// business errors
var ErrUserAlreadyExists = errors.New("user with that email already exists")
