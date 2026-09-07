package models

import "errors"

var (
	ErrNotFound              = errors.New("resource not found")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrInsufficientInventory = errors.New("insufficient inventory")
	ErrBadRequest            = errors.New("bad request")
	ErrConflict              = errors.New("resource conflict")
	ErrInternal              = errors.New("internal server error")
	ErrValidationError       = errors.New("validation failed")
)
