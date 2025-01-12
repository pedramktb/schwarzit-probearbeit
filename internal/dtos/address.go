package dtos

import "github.com/pedramktb/schwarzit-probearbeit/internal/types"

// @Description Address DTO model for retrievals, creations and updates
type Address struct {
	Street       string  `json:"street" binding:"required" validate:"required" example:"Main Street"`
	StreetNumber string  `json:"street_number" binding:"required" validate:"required" example:"123"`
	Extra        *string `json:"extra" binding:"omitnil" validate:"optional" example:"Apartment 1"`
	ZipCode      string  `json:"zip_code" binding:"required" validate:"required" example:"12345"`
	City         string  `json:"city" binding:"required" validate:"required" example:"Berlin"`
} // @name Address

func FromAddress(c *types.Address) Address {
	return Address{
		Street:       c.Street,
		StreetNumber: c.StreetNumber,
		Extra:        c.Extra,
		ZipCode:      c.ZipCode,
		City:         c.City,
	}
}

func (c *Address) ToAddress() types.Address {
	return types.Address{
		Street:       c.Street,
		StreetNumber: c.StreetNumber,
		Extra:        c.Extra,
		ZipCode:      c.ZipCode,
		City:         c.City,
	}
}
