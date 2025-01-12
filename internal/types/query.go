package types

import (
	"gorm.io/gorm"
)

type Mappable interface {
	ToMap() map[string]any
}

type Pagination struct {
	Offset int
	Limit  int
}

type QueryParams struct {
	Conditions Mappable
	Pagination
}

const defaultLimit = 10

func Query(tx *gorm.DB, params QueryParams) *gorm.DB {
	if params.Limit == 0 {
		params.Limit = defaultLimit
	}

	if params.Conditions == nil {
		return tx.Limit(params.Limit).Offset(params.Offset)
	} else {
		return tx.Where(params.Conditions.ToMap()).Limit(params.Limit).Offset(params.Offset)
	}
}
