package models

import "github.com/google/uuid"

type Team struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID        uuid.UUID  `gorm:"type:uuid;index"`
	ParentTeamID *uuid.UUID `gorm:"type:uuid;index"`
	Name         string
	Path         string
	Timestamps
	SoftDelete
}

func (Team) TableName() string { return "teams" }
