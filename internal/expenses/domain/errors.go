package domain

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrInvalidInput          = errors.New("invalid input")
	ErrConflict              = errors.New("conflict")
	ErrForbidden             = errors.New("forbidden")
	ErrInsufficientBudget    = errors.New("insufficient budget")
	ErrPolicyViolation       = errors.New("policy violation")
	ErrExpenseNotEditable    = errors.New("expense not editable in current status")
	ErrDuplicateReceipt      = errors.New("duplicate receipt detected")
	ErrAdvanceAlreadySettled = errors.New("advance already settled")
	ErrReimbursementPaid     = errors.New("reimbursement already paid")
)
