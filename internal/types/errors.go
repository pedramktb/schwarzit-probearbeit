package types

import (
	"errors"

	"gorm.io/gorm"
)

// Errors
var (
	// Primary Errors
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInternal     = errors.New("internal error")

	// ErrBadRequest Most Used Secondary Errors
	ErrInvalidID = errors.Join(ErrBadRequest, errors.New("invalid id"))

	// ErrInternal Most Used Secondary Errors
	ErrDBUnhandled   = errors.Join(ErrInternal, errors.New("database unhandled error"))
	ErrDataImmutable = errors.Join(ErrInternal, errors.New("data is immutable"))
	ErrDataCorrupted = errors.Join(ErrInternal, errors.New("data corrupted"))
)

// DBError converts gorm errors to internal errors
func DBError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrCheckConstraintViolated),
		errors.Is(err, gorm.ErrForeignKeyViolated):
		return errors.Join(ErrBadRequest, err)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errors.Join(ErrNotFound, err)
	case errors.Is(err, gorm.ErrDryRunModeUnsupported),
		errors.Is(err, gorm.ErrEmptySlice),
		errors.Is(err, gorm.ErrInvalidDB),
		errors.Is(err, gorm.ErrDuplicatedKey),
		errors.Is(err, gorm.ErrInvalidData),
		errors.Is(err, gorm.ErrInvalidField),
		errors.Is(err, gorm.ErrInvalidTransaction),
		errors.Is(err, gorm.ErrInvalidValue),
		errors.Is(err, gorm.ErrInvalidValueOfLength),
		errors.Is(err, gorm.ErrMissingWhereClause),
		errors.Is(err, gorm.ErrModelAccessibleFieldsRequired),
		errors.Is(err, gorm.ErrModelValueRequired),
		errors.Is(err, gorm.ErrNotImplemented),
		errors.Is(err, gorm.ErrPreloadNotAllowed),
		errors.Is(err, gorm.ErrPrimaryKeyRequired),
		errors.Is(err, gorm.ErrRegistered),
		errors.Is(err, gorm.ErrSubQueryRequired),
		errors.Is(err, gorm.ErrUnsupportedDriver),
		errors.Is(err, gorm.ErrUnsupportedRelation):
		return errors.Join(ErrInternal, err)
	}
	return ErrDBUnhandled
}
