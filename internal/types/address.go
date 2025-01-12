package types

import (
	"database/sql/driver"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
)

type Address struct {
	Street       string
	StreetNumber string
	Extra        *string
	ZipCode      string
	City         string
}

func (a *Address) Scan(value any) error {
	var s string
	switch v := value.(type) {
	case nil:
		return nil
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("unsupported type for Address: %T", value)
	}

	s = strings.TrimSuffix(strings.TrimPrefix(s, "("), ")")
	reader := csv.NewReader(strings.NewReader(s))
	reader.TrimLeadingSpace = true
	fields, err := reader.Read()
	if err != nil {
		return errors.Newf("error reading CSV: %v", err)
	}

	if len(fields) != 5 {
		return errors.Newf("expected 4 fields, got %d", len(fields))
	}

	a.Street = fields[0]
	a.StreetNumber = fields[1]
	if fields[2] != "" && fields[2] != "NULL" {
		a.Extra = &fields[2]
	}
	a.ZipCode = fields[3]
	a.City = fields[4]

	return nil
}

func (a Address) Value() (driver.Value, error) {
	if a.Extra != nil {
		return fmt.Sprintf("(%q,%q,%q,%q,%q)", a.Street, a.StreetNumber, *a.Extra, a.ZipCode, a.City), nil
	} else {
		return fmt.Sprintf("(%q,%q,NULL,%q,%q)", a.Street, a.StreetNumber, a.ZipCode, a.City), nil
	}
}

type AddressPatch struct {
	Street       Optional[string]
	StreetNumber Optional[string]
	Extra        Optional[*string]
	ZipCode      Optional[string]
	City         Optional[string]
}

func (c *AddressPatch) ToMap() map[string]any {
	m := make(map[string]any)
	if c.Street.HasValue {
		m["street"] = c.Street.Value
	}
	if c.StreetNumber.HasValue {
		m["street_number"] = c.StreetNumber.Value
	}
	if c.Extra.HasValue {
		m["extra"] = c.Extra.Value
	}
	if c.ZipCode.HasValue {
		m["zip_code"] = c.ZipCode.Value
	}
	if c.City.HasValue {
		m["city"] = c.City.Value
	}
	return m
}
