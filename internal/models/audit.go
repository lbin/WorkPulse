package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditLog struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID       uuid.UUID      `gorm:"type:uuid;index"`
	ActorUserID uuid.UUID      `gorm:"type:uuid;index;column:actor_user_id"`
	Action      string         `gorm:"column:action"`
	EntityType  string         `gorm:"column:entity_type"`
	EntityID    uuid.UUID      `gorm:"type:uuid;column:entity_id"`
	Diff        datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
}

func (AuditLog) TableName() string { return "audit_logs" }
