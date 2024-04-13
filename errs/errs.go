package errs

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrServerErr      = errors.New("server processing error")
)
