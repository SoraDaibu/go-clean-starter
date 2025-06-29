package user

import "errors"

var (
	ErrNameIsRequired     = errors.New("name is required")
	ErrEmailIsRequired    = errors.New("email is required")
	ErrPasswordIsRequired = errors.New("password is required")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters long")
)
