package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WorkItem struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID        uuid.UUID  `gorm:"type:uuid;index"`
	TeamID       *uuid.UUID `gorm:"type:uuid;index"`
	OwnerUserID  uuid.UUID  `gorm:"type:uuid;index;column:owner_user_id"`
	Type         string
	Title        string
	Description  *string
	Status       string
	Priority     int
	EffortMinutes *int `gorm:"column:effort_minutes"`
	StartAt      *time.Time `gorm:"column:start_at"`
	EndAt        *time.Time `gorm:"column:end_at"`
	DueAt        *time.Time `gorm:"column:due_at"`
	SourceType   *string    `gorm:"column:source_type"`
	SourceID     *uuid.UUID `gorm:"type:uuid;column:source_id"`
	ReasonCode   *string    `gorm:"column:reason_code"`
	Tags         datatypes.JSON `gorm:"type:jsonb"`
	SchemaVersion int `gorm:"column:schema_version"`
	CustomFields datatypes.JSON `gorm:"type:jsonb;column:custom_fields"`
	Timestamps
	SoftDelete
}

func (WorkItem) TableName() string { return "work_items" }

type WorkItemOKRLink struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID              uuid.UUID  `gorm:"type:uuid;index"`
	WorkItemID         uuid.UUID  `gorm:"type:uuid;index;column:work_item_id"`
	ObjectiveID        *uuid.UUID `gorm:"type:uuid;index;column:objective_id"`
	KeyResultID        *uuid.UUID `gorm:"type:uuid;index;column:key_result_id"`
	LinkType           string     `gorm:"column:link_type"`
	ContributionWeight float64    `gorm:"column:contribution_weight"`
	EvidenceRequired   bool       `gorm:"column:evidence_required"`
	CreatedAt          time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (WorkItemOKRLink) TableName() string { return "work_item_okr_links" }
