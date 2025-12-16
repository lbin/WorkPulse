package models

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID        uuid.UUID `gorm:"type:uuid;index"`
	Email        string
	PasswordHash string
	DisplayName  string
	AvatarURL    *string
	Status       string
	Timestamps
	SoftDelete
}

func (User) TableName() string { return "users" }
