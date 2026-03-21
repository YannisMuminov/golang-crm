package apperror

import (
	"errors"
	"net/http"
)

var (
	ErrEmailTaken      = errors.New("email already taken")
	ErrInvalidCreds    = errors.New("invalid credentials")
	ErrAccountDisabled = errors.New("account disabled")
	ErrUnauthorized    = errors.New("unauthorized")

	ErrNotFound   = errors.New("not found")
	ErrForbidden  = errors.New("forbidden")
	ErrBadRequest = errors.New("bad request")
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func HTTPStatus(err error) (int, string) {
	var appErr *AppError

	if errors.As(err, &appErr) {
		return appErr.Code, appErr.Message
	}

	switch {
	case errors.Is(err, ErrEmailTaken):
		return http.StatusConflict, err.Error()
	case errors.Is(err, ErrInvalidCreds):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, ErrAccountDisabled):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
