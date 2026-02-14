package apperror

import "fmt"

const (
	CodeNotFound            = "NOT_FOUND"
	CodeConflict            = "CONFLICT"
	CodeInsufficientBalance = "INSUFFICIENT_BALANCE"
	CodeValidation          = "VALIDATION_ERROR"
	CodeInternal            = "INTERNAL_ERROR"
)

type AppError interface {
	error
	Code() string
}

type ErrNotFound struct {
	Entity string
	ID     int64
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found. Please check the ID and try again", e.Entity)
}

func (e *ErrNotFound) Code() string { return CodeNotFound }

type ErrConflict struct {
	Entity string
	ID     int64
}

func (e *ErrConflict) Error() string {
	return fmt.Sprintf("This %s already exists in the system", e.Entity)
}

func (e *ErrConflict) Code() string { return CodeConflict }

type ErrInsufficientBalance struct {
	AccountID int64
}

func (e *ErrInsufficientBalance) Error() string {
	return "Insufficient funds. Please check your account balance and try again"
}

func (e *ErrInsufficientBalance) Code() string { return CodeInsufficientBalance }

type ErrValidation struct {
	Message string
}

func (e *ErrValidation) Error() string {
	return e.Message
}

func (e *ErrValidation) Code() string { return CodeValidation }
