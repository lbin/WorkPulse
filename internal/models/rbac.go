package models

import (
	"time"

	"github.com/google/uuid"
)

type OrgUnit struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID        uuid.UUID  `gorm:"type:uuid;index"`
	ParentUnitID *uuid.UUID `gorm:"type:uuid;column:parent_unit_id"`
	Name         string
	Path         string
	UnitType     string
	Timestamps
	SoftDelete
}

func (OrgUnit) TableName() string { return "org_units" }

type Role struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID       uuid.UUID  `gorm:"type:uuid;index"`
	OrgUnitID   *uuid.UUID `gorm:"type:uuid;column:org_unit_id"`
	Name        string
	Description *string
	Timestamps
	SoftDelete
}

func (Role) TableName() string { return "roles" }

type Permission struct {
	Code        string    `gorm:"primaryKey;column:code"`
	Description *string   `gorm:"column:description"`
	Scope       *string   `gorm:"column:scope"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Permission) TableName() string { return "permissions" }

type RolePermission struct {
	RoleID         uuid.UUID `gorm:"type:uuid;primaryKey;column:role_id"`
	PermissionCode string    `gorm:"primaryKey;column:permission_code"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (RolePermission) TableName() string { return "role_permissions" }

type RoleBinding struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID      uuid.UUID `gorm:"type:uuid;index"`
	RoleID     uuid.UUID `gorm:"type:uuid;column:role_id"`
	EntityType string    `gorm:"column:entity_type"`
	EntityID   uuid.UUID `gorm:"type:uuid;column:entity_id"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (RoleBinding) TableName() string { return "role_bindings" }

type Membership struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID      uuid.UUID  `gorm:"type:uuid;index"`
	TeamID     *uuid.UUID `gorm:"type:uuid;column:team_id"`
	OrgUnitID  *uuid.UUID `gorm:"type:uuid;column:org_unit_id"`
	UserID     uuid.UUID  `gorm:"type:uuid;column:user_id"`
	RoleID     *uuid.UUID `gorm:"type:uuid;column:role_id"`
	RoleInTeam string     `gorm:"column:role_in_team"`
	Timestamps
	SoftDelete
}

func (Membership) TableName() string { return "memberships" }
