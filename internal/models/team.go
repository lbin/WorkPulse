package models

import "github.com/google/uuid"

type Team struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrgID        uuid.UUID  `gorm:"type:uuid;index" json:"org_id"`
	ParentTeamID *uuid.UUID `gorm:"type:uuid;index" json:"parent_team_id"`
	Name         string     `json:"name"`
	Path         string     `json:"path"`
	Timestamps
	SoftDelete
}

func (Team) TableName() string { return "teams" }
