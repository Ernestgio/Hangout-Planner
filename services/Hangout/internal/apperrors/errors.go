package apperrors

import "errors"

var ErrUserAlreadyExists = errors.New("user with that email already exists")
var ErrInvalidPayload = errors.New("invalid payload")
