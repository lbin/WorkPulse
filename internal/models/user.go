package models

import "github.com/google/uuid"

// User represents an account within an organization.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrgID        uuid.UUID `gorm:"type:uuid;index" json:"org_id"`
	PasswordHash string    `gorm:"column:password_hash" json:"-"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	AvatarURL    *string   `json:"avatar_url"`
	Status       string    `json:"status"`
	Timestamps
	SoftDelete
}

func (User) TableName() string { return "users" }
