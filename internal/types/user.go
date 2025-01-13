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

func (u *User) ToSave() (base, version map[string]any) {
	u.VersionID = uuid.New()
	u.IsLatestVersion = true

	if u.ID == uuid.Nil {
		u.ID = uuid.New()
		base = map[string]any{
			"id": u.ID,
		}
	}

	version = map[string]any{
		"id":            u.VersionID,
		"user_id":       u.ID,
		"first_name":    u.FirstName,
		"last_name":     u.LastName,
		"email":         u.Email,
		"phone":         u.Phone,
		"is_admin":      u.IsAdmin,
		"password_hash": u.PasswordHash,
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

func (u *UserPatch) ToMap() map[string]any {
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
