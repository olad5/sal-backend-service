package errors

import "errors"

const (
	ErrSomethingWentWrong = "something went wrong"
	ErrInvalidJson        = "Invalid JSON"
	ErrMissingBody        = "missing body request"
	ErrUnauthorized       = "unauthorized to perform this action"
)

var ErrInvalidID = errors.New("ID is not in its proper form")
