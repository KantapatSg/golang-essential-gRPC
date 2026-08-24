package service

import "errors"

var (
	ErrNotFound      = errors.New("product not found")
	ErrInvalidInput  = errors.New("invalid product input")
	ErrAlreadyExists = errors.New("product already exists")
)
