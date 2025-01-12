package types

import "github.com/google/uuid"

type User struct {
	ID              uuid.UUID `gorm:"column:id"`
	VersionID       uuid.UUID `gorm:"column:version_id"`
	IsLatestVersion bool      `gorm:"column:is_latest_version"`
	FirstName       string    `gorm:"column:first_name"`
	LastName        string    `gorm:"column:last_name"`
	Email           string    `gorm:"column:email"`
	Phone           string    `gorm:"column:phone"`
	IsAdmin         bool      `gorm:"column:is_admin"`
	PasswordHash    string    `gorm:"column:password_hash"`
}

func (c *User) ToSave() (base, version map[string]any) {
	c.VersionID = uuid.New()
	c.IsLatestVersion = true

	if c.ID == uuid.Nil {
		c.ID = uuid.New()
		base = map[string]any{
			"id": c.ID,
		}
	}

	version = map[string]any{
		"id":            c.VersionID,
		"user_id":       c.ID,
		"first_name":    c.FirstName,
		"last_name":     c.LastName,
		"email":         c.Email,
		"phone":         c.Phone,
		"is_admin":      c.IsAdmin,
		"password_hash": c.PasswordHash,
	}

	return base, version
}

type UserPatch struct {
	ID           Optional[uuid.UUID]
	VersionID    Optional[uuid.UUID]
	FirstName    Optional[string]
	LastName     Optional[string]
	Email        Optional[string]
	Phone        Optional[string]
	IsAdmin      Optional[bool]
	PasswordHash Optional[string]
}

func (u UserPatch) ToMap() map[string]any {
	m := make(map[string]any)
	if u.ID.HasValue {
		m["id"] = u.ID.Value
	}
	if u.VersionID.HasValue {
		m["version_id"] = u.VersionID.Value
	}
	if u.FirstName.HasValue {
		m["first_name"] = u.FirstName.Value
	}
	if u.LastName.HasValue {
		m["last_name"] = u.LastName.Value
	}
	if u.Email.HasValue {
		m["email"] = u.Email.Value
	}
	if u.Phone.HasValue {
		m["phone"] = u.Phone.Value
	}
	if u.IsAdmin.HasValue {
		m["is_admin"] = u.IsAdmin.Value
	}
	if u.PasswordHash.HasValue {
		m["password_hash"] = u.PasswordHash.Value
	}
	return m
}

func (u *User) ApplyPatch(p UserPatch) {
	if p.FirstName.HasValue {
		u.FirstName = p.FirstName.Value
	}
	if p.LastName.HasValue {
		u.LastName = p.LastName.Value
	}
	if p.Email.HasValue {
		u.Email = p.Email.Value
	}
	if p.Phone.HasValue {
		u.Phone = p.Phone.Value
	}
	if p.IsAdmin.HasValue {
		u.IsAdmin = p.IsAdmin.Value
	}
	if p.PasswordHash.HasValue {
		u.PasswordHash = p.PasswordHash.Value
	}
}
