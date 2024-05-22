package domain

import "errors"

var (
	ErrInvalidParams = errors.New("invalid body")
	ErrNotFound      = errors.New("mark is not found")
)
