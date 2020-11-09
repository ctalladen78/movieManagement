package errs

import "errors"

var (
	// ErrNotFound is a not found error
	ErrNotFound = errors.New("not found")
	// ErrInvalid is an invalid request error
	ErrInvalid = errors.New("invalid request")
	// ErrUnauthorized is an unauthorized error
	ErrUnauthorized = errors.New("unauthorized")
	// ErrConflict is an conflict/duplicate error
	ErrConflict = errors.New("conflict")
	// ErrForbidden is an forbidden error
	ErrForbidden = errors.New("forbidden")
)
