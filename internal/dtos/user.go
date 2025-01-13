package dtos

import (
	"github.com/google/uuid"
	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
	"golang.org/x/crypto/bcrypt"
)

// @Description User DTO model for responses
// @Tags user
type User struct {
	ID        uuid.UUID `json:"id" swaggertype:"string" format:"uuid" example:"b05a5d28-1a51-46a8-b35c-6e160a05a0ad"`
	VersionID uuid.UUID `json:"version_id" swaggertype:"string" format:"uuid" example:"b05a5d28-1a51-46a8-b35c-6e160a05a0ad"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	Email     string    `json:"email" format:"email" example:"abc@xyz.com"`
	Phone     string    `json:"phone" format:"phone" example:"+49123456789"`
} // @name User

// @Description QueryUser DTO model for user queries
// @Tags user
type QueryUser struct {
	ID        *uuid.UUID `json:"id" swaggertype:"string" format:"uuid" example:"b05a5d28-1a51-46a8-b35c-6e160a05a0ad"`
	VersionID *uuid.UUID `json:"version_id" swaggertype:"string" format:"uuid" example:"b05a5d28-1a51-46a8-b35c-6e160a05a0ad"`
	FirstName *string    `json:"first_name" example:"John"`
	LastName  *string    `json:"last_name" example:"Doe"`
	Email     *string    `json:"email" binding:"omitempty,email" format:"email" example:"abc@xyz.com"`
	Phone     *string    `json:"phone" format:"phone" example:"+49123456789"`
} // @name QueryUser

// @Description UserQueryParams DTO model for user query parameters
// @Tags user
type UserQueryParams struct {
	Pagination
	QueryUser
} // @name UserQueryParams

// @Description SaveUser DTO model for user creation and updates (overwrites)
// @Tags user
type SaveUser struct {
	FirstName string `json:"first_name" binding:"required" validate:"required" example:"John"`
	LastName  string `json:"last_name" binding:"required" validate:"required" example:"Doe"`
	Email     string `json:"email" binding:"required,email" validate:"required" format:"email" example:"abc@xyz.com"`
	Phone     string `json:"phone" binding:"required" validate:"required" format:"phone" example:"+49123456789"`
	IsAdmin   bool   `json:"is_admin" example:"false"`
	Password  string `json:"password" binding:"required" validate:"required" example:"password"`
} // @name SaveUser

// @Description PatchUser DTO model for user updates (partial)
// @Tags user
type PatchUser struct {
	FirstName *string `json:"first_name" example:"John"`
	LastName  *string `json:"last_name" example:"Doe"`
	Email     *string `json:"email" binding:"omitempty,email" format:"email" example:"abc@xyz.com"`
	Phone     *string `json:"phone" format:"phone" example:"+49123456789"`
	IsAdmin   *bool   `json:"is_admin" example:"false"`
	Password  *string `json:"password" example:"password"`
} // @name PatchUser

// @Description RegisterUser DTO model for user registration
// @Tags user
type RegisterUser struct {
	FirstName string `json:"first_name" binding:"required" validate:"required" example:"John"`
	LastName  string `json:"last_name" binding:"required" validate:"required" example:"Doe"`
	Email     string `json:"email" binding:"required,email" validate:"required" format:"email" example:"abc@xyz.com"`
	Phone     string `json:"phone" binding:"required" validate:"required" example:"+49123456789"`
	Password  string `json:"password" binding:"required" validate:"required" example:"password"`
} // @name RegisterUser

func FromUser(u *types.User) User {
	return User{
		ID:        u.ID,
		VersionID: u.VersionID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
	}
}

func (u *QueryUser) ToUserPatch() types.UserPatch {
	p := types.UserPatch{}
	if u.ID != nil {
		p.ID = types.Optional[uuid.UUID]{HasValue: true, Value: *u.ID}
	}
	if u.VersionID != nil {
		p.VersionID = types.Optional[uuid.UUID]{HasValue: true, Value: *u.VersionID}
	}
	if u.FirstName != nil {
		p.FirstName = types.Optional[string]{HasValue: true, Value: *u.FirstName}
	}
	if u.LastName != nil {
		p.LastName = types.Optional[string]{HasValue: true, Value: *u.LastName}
	}
	if u.Email != nil {
		p.Email = types.Optional[string]{HasValue: true, Value: *u.Email}
	}
	if u.Phone != nil {
		p.Phone = types.Optional[string]{HasValue: true, Value: *u.Phone}
	}
	return p
}

func (u *UserQueryParams) ToQueryParams() types.QueryParams {
	return types.QueryParams{
		Pagination: u.Pagination.ToPagination(),
		Conditions: types.Pointer(u.QueryUser.ToUserPatch()),
	}
}

func (u *SaveUser) ToCreateUser() types.User {
	return types.User{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		Phone:        u.Phone,
		IsAdmin:      u.IsAdmin,
		PasswordHash: HashPassword(u.Password),
	}
}

func (u *SaveUser) ToUpdateUser(id uuid.UUID) types.User {
	return types.User{
		ID:           id,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		Phone:        u.Phone,
		IsAdmin:      u.IsAdmin,
		PasswordHash: HashPassword(u.Password),
	}
}

func (u *PatchUser) ToUserPatch() types.UserPatch {
	userPatch := types.UserPatch{}
	if u.FirstName != nil {
		userPatch.FirstName = types.Optional[string]{HasValue: true, Value: *u.FirstName}
	}
	if u.LastName != nil {
		userPatch.LastName = types.Optional[string]{HasValue: true, Value: *u.LastName}
	}
	if u.Email != nil {
		userPatch.Email = types.Optional[string]{HasValue: true, Value: *u.Email}
	}
	if u.Phone != nil {
		userPatch.Phone = types.Optional[string]{HasValue: true, Value: *u.Phone}
	}
	if u.IsAdmin != nil {
		userPatch.IsAdmin = types.Optional[bool]{HasValue: true, Value: *u.IsAdmin}
	}
	if u.Password != nil {
		userPatch.PasswordHash = types.Optional[string]{HasValue: true, Value: HashPassword(*u.Password)}
	}
	return userPatch
}

func (u *RegisterUser) ToUser() types.User {
	return types.User{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		Phone:        u.Phone,
		PasswordHash: HashPassword(u.Password),
	}
}

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
