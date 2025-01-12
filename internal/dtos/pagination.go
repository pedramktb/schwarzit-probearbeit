package dtos

import "github.com/pedramktb/schwarzit-probearbeit/internal/types"

// @Description pagination model for queries
// @Tags pagination
type Pagination struct {
	Limit  int `json:"limit" example:"10"`
	Offset int `json:"offset" example:"0"`
} // @name Pagination

func (p *Pagination) ToPagination() types.Pagination {
	return types.Pagination{
		Limit:  p.Limit,
		Offset: p.Offset,
	}
}
