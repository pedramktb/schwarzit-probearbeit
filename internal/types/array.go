package types

import (
	"database/sql/driver"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
)

// Array is a custom type that implements the sql.Scanner and driver.Valuer interfaces.
// This is important for the gorm library to be able to store arrays in the database.
// Example: gorm fails to store []string but can store Array[string]
// Note: only primitive types are supported as complex types can have different
// representations in string or byte slice format. See the time_slots.go as an example.

type Array[T string |
	bool |
	uintptr |
	int | int64 | int32 | int16 | int8 |
	uint | uint64 | uint32 | uint16 | uint8 |
	float64 | float32 |
	complex128 | complex64] []T

func (a *Array[T]) Scan(value any) error {
	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return errors.Newf("unsupported type for Array: %T", value)
	}

	*a = []T{}

	s = strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
	parts := strings.Split(s, ",")
	for _, part := range parts {
		t, err := stringToType[T](part)
		if err != nil {
			return err
		}

		*a = append(*a, t)
	}

	return nil
}

func (a Array[T]) Value() (driver.Value, error) {
	return []T(a), nil
}

//nolint:gocyclo // this function is long but it's just a switch case
func stringToType[T string |
	bool |
	uintptr |
	int | int64 | int32 | int16 | int8 |
	uint | uint64 | uint32 | uint16 | uint8 |
	float64 | float32 |
	complex128 | complex64](s string) (T, error) {
	var t T
	var val any
	var err error

	var i64 int64
	var u64 uint64
	var f64 float64
	var c128 complex128
	switch any(t).(type) {
	case string:
		val = s
	case bool:
		val, err = strconv.ParseBool(s)
	case uintptr:
		u64, err = strconv.ParseUint(s, 10, 0)
		val = uintptr(u64)
	case int:
		i64, err = strconv.ParseInt(s, 10, 0)
		val = int(i64)
	case int64:
		val, err = strconv.ParseInt(s, 10, 64)
	case int32:
		i64, err = strconv.ParseInt(s, 10, 32)
		val = int32(i64)
	case int16:
		i64, err = strconv.ParseInt(s, 10, 16)
		val = int16(i64)
	case int8:
		i64, err = strconv.ParseInt(s, 10, 8)
		val = int8(i64)
	case uint:
		u64, err = strconv.ParseUint(s, 10, 0)
		val = uint(u64)
	case uint64:
		val, err = strconv.ParseUint(s, 10, 64)
	case uint32:
		u64, err = strconv.ParseUint(s, 10, 32)
		val = uint32(u64)
	case uint16:
		u64, err = strconv.ParseUint(s, 10, 16)
		val = uint16(u64)
	case uint8:
		u64, err = strconv.ParseUint(s, 10, 8)
		val = uint8(u64)
	case float64:
		val, err = strconv.ParseFloat(s, 64)
	case float32:
		f64, err = strconv.ParseFloat(s, 32)
		val = float32(f64)
	case complex128:
		val, err = strconv.ParseComplex(s, 128)
	case complex64:
		c128, err = strconv.ParseComplex(s, 64)
		val = complex64(c128)
	default:
		return t, errors.Newf("unsupported type: %T", t)
	}
	t = val.(T)
	return t, err
}
