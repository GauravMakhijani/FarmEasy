package services

import "errors"

var (
	ErrUnauthorized   = errors.New("incorrect email or password")
	ErrDuplicateEmail = errors.New("account exists for the given email")
	ErrDuplicatePhone = errors.New("account exists for the given phone")
)
