package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OKRCycle struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID         uuid.UUID  `gorm:"type:uuid;index"`
	TeamID        *uuid.UUID `gorm:"type:uuid;index"`
	Type          string
	Name          string
	StartDate     time.Time `gorm:"type:date;column:start_date"`
	EndDate       time.Time `gorm:"type:date;column:end_date"`
	Status        string
	CreatedBy     uuid.UUID `gorm:"type:uuid;column:created_by"`
	SchemaVersion int       `gorm:"column:schema_version"`
	Timestamps
	SoftDelete
}

func (OKRCycle) TableName() string { return "okr_cycles" }

type OKRObjective struct {
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

func (OKRObjective) TableName() string { return "okr_objectives" }

type OKRKeyResult struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID         uuid.UUID `gorm:"type:uuid;index"`
	ObjectiveID   uuid.UUID `gorm:"type:uuid;index;column:objective_id"`
	Title         string
	MetricType    string   `gorm:"column:metric_type"`
	TargetValue   *float64 `gorm:"column:target_value"`
	CurrentValue  *float64 `gorm:"column:current_value"`
	Unit          *string
	Weight        float64
	Confidence    int
	Status        string
	SchemaVersion int            `gorm:"column:schema_version"`
	MetricPayload datatypes.JSON `gorm:"type:jsonb;column:metric_payload"`
	Timestamps
	SoftDelete
}

func (OKRKeyResult) TableName() string { return "okr_key_results" }

type OKRLink struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID         uuid.UUID  `gorm:"type:uuid;index"`
	ObjectiveID   *uuid.UUID `gorm:"type:uuid;index;column:objective_id"`
	KeyResultID   *uuid.UUID `gorm:"type:uuid;index;column:key_result_id"`
	EntityType    string
	EntityID      uuid.UUID `gorm:"type:uuid;column:entity_id"`
	Relation      string
	SchemaVersion int       `gorm:"column:schema_version"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (OKRLink) TableName() string { return "okr_links" }
