package handlers

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func NewApiError() *ApiError {
	return &ApiError{}
}

func (e *ApiError) Error(status int, message string, err error) {
	var pgError *pgconn.PgError

	e.Status = status
	e.Message = message

	if errors.As(err, &pgError) {
		e.Details = pgError.Message
	}

	if !errors.As(err, &pgError) {
		e.Details = err.Error()
	}
}

func ApiErrorResponse(status int, message string, err error) *ApiError {
	newError := NewApiError()
	newError.Error(status, message, err)
	return newError
}

func HandleError(status int, message string, err error) *ApiError {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch pgError.Code {
		case "23505":
			return ApiErrorResponse(409, message, pgError)
		case "23503":
			return ApiErrorResponse(400, message, pgError)
		case "23514":
			return ApiErrorResponse(409, message, pgError)
		default:
			return ApiErrorResponse(500, message, pgError)
		}
	}
	return ApiErrorResponse(status, message, err)

}
