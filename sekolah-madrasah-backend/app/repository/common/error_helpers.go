package common

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

// Repository error constants
const (
	ErrIDRequired          = "id is required"
	ErrRecordNotFound      = "record not found"
	ErrInvalidInput        = "invalid input provided"
	ErrTypeAssertionFailed = "type assertion failed"
)

// HandleGORMError standardizes GORM error handling across repositories
func HandleGORMError(err error) (int, error) {
	if err == nil {
		return http.StatusOK, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, err
	}

	return http.StatusInternalServerError, err
}

// ValidateRequiredID checks if required ID is provided
func ValidateRequiredID(id interface{}) (int, error) {
	if id == nil {
		return http.StatusBadRequest, errors.New(ErrIDRequired)
	}
	return http.StatusOK, nil
}
