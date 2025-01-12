package types

import (
	"github.com/cockroachdb/errors"
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
	ErrInvalidID = errors.Wrap(ErrBadRequest, "invalid id")

	// ErrInternal Most Used Secondary Errors
	ErrDBUnhandled   = errors.Wrap(ErrInternal, "database unhandled error")
	ErrDataImmutable = errors.Wrap(ErrInternal, "data is immutable")
	ErrDataCorrupted = errors.Wrap(ErrInternal, "data corrupted")
)

// DBError converts gorm errors to internal errors
func DBError(err error) error {
	if err == nil {
		return nil
	} else if errors.IsAny(err,
		gorm.ErrCheckConstraintViolated,
		gorm.ErrForeignKeyViolated,
	) {
		return errors.CombineErrors(ErrBadRequest, err)
	} else if errors.IsAny(err,
		gorm.ErrRecordNotFound,
	) {
		return errors.CombineErrors(ErrNotFound, err)
	} else if errors.IsAny(err,
		gorm.ErrDryRunModeUnsupported,
		gorm.ErrDuplicatedKey,
		gorm.ErrEmptySlice,
		gorm.ErrInvalidDB,
		gorm.ErrDuplicatedKey,
		gorm.ErrInvalidData,
		gorm.ErrInvalidField,
		gorm.ErrInvalidTransaction,
		gorm.ErrInvalidValue,
		gorm.ErrInvalidValueOfLength,
		gorm.ErrMissingWhereClause,
		gorm.ErrModelAccessibleFieldsRequired,
		gorm.ErrModelValueRequired,
		gorm.ErrNotImplemented,
		gorm.ErrPreloadNotAllowed,
		gorm.ErrPrimaryKeyRequired,
		gorm.ErrRegistered,
		gorm.ErrSubQueryRequired,
		gorm.ErrUnsupportedDriver,
		gorm.ErrUnsupportedRelation,
	) {
		return errors.CombineErrors(ErrInternal, err)
	}
	return ErrDBUnhandled
}
