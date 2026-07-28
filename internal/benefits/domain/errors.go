package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrConflict        = errors.New("conflict")
	ErrForbidden       = errors.New("forbidden")
	ErrIneligible      = errors.New("employee not eligible")
	ErrWorkflow        = errors.New("workflow error")
	ErrBalanceExceeded = errors.New("balance exceeded")
)
