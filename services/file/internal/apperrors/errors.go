package apperrors

import "errors"

// application errors
var ErrAppPortRequired = errors.New("APP_PORT is required")
var ErrDbPasswordRequired = errors.New("DB_PASSWORD required in production")
