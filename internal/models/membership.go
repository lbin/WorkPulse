package models

import "github.com/google/uuid"

// Membership connects a user to a team within an organization.
type Membership struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrgID      uuid.UUID `gorm:"type:uuid;index" json:"org_id"`
	TeamID     uuid.UUID `gorm:"type:uuid;index" json:"team_id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	RoleInTeam string    `json:"role_in_team"`
	Timestamps
	SoftDelete
}

// TableName keeps the mapping aligned with the database table.
func (Membership) TableName() string { return "memberships" }
