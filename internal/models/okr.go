package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OKRCycle struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID     uuid.UUID  `gorm:"type:uuid;index"`
	TeamID    *uuid.UUID `gorm:"type:uuid;index"`
	Type      string
	Name      string
	StartDate time.Time `gorm:"type:date;column:start_date"`
	EndDate   time.Time `gorm:"type:date;column:end_date"`
	Status    string
	CreatedBy uuid.UUID `gorm:"type:uuid;column:created_by"`
	Timestamps
	SoftDelete
}

func (OKRCycle) TableName() string { return "okr_cycles" }

type Objective struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID             uuid.UUID  `gorm:"type:uuid;index"`
	CycleID           uuid.UUID  `gorm:"type:uuid;index;column:cycle_id"`
	TeamID            *uuid.UUID `gorm:"type:uuid;index"`
	OwnerUserID       uuid.UUID  `gorm:"type:uuid;index;column:owner_user_id"`
	ParentObjectiveID *uuid.UUID `gorm:"type:uuid;index;column:parent_objective_id"`
	Title             string
	Description       *string
	Weight            float64
	Status            string
	Tags              datatypes.JSON `gorm:"type:jsonb"`
	SchemaVersion     int            `gorm:"column:schema_version"`
	Payload           datatypes.JSON `gorm:"type:jsonb"`
	Timestamps
	SoftDelete
}

func (Objective) TableName() string { return "objectives" }

type KeyResult struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID        uuid.UUID `gorm:"type:uuid;index"`
	ObjectiveID  uuid.UUID `gorm:"type:uuid;index;column:objective_id"`
	Title        string
	MetricType   string `gorm:"column:metric_type"`
	TargetValue  *float64 `gorm:"column:target_value"`
	CurrentValue *float64 `gorm:"column:current_value"`
	Unit         *string
	Weight       float64
	Confidence   int
	Status       string
	SchemaVersion int `gorm:"column:schema_version"`
	MetricPayload datatypes.JSON `gorm:"type:jsonb;column:metric_payload"`
	Timestamps
	SoftDelete
}

func (KeyResult) TableName() string { return "key_results" }
