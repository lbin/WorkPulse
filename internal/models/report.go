package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Report struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID         uuid.UUID  `gorm:"type:uuid;index"`
	TeamID        *uuid.UUID `gorm:"type:uuid;index"`
	AuthorUserID  uuid.UUID  `gorm:"type:uuid;index;column:author_user_id"`
	Type          string
	PeriodStart   time.Time `gorm:"type:date;column:period_start"`
	PeriodEnd     time.Time `gorm:"type:date;column:period_end"`
	Status        string
	Title         *string
	Summary       *string
	SchemaVersion int            `gorm:"column:schema_version"`
	Payload       datatypes.JSON `gorm:"type:jsonb"`
	Timestamps
	SoftDelete
}

func (Report) TableName() string { return "reports" }

type ReportEntry struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID        uuid.UUID  `gorm:"type:uuid;index"`
	ReportID     uuid.UUID  `gorm:"type:uuid;index;column:report_id"`
	Section      string
	Content      string
	EffortMinutes *int `gorm:"column:effort_minutes"`
	Status       string
	WorkItemID   *uuid.UUID `gorm:"type:uuid;index;column:work_item_id"`
	OrderNo      int        `gorm:"column:order_no"`
	Timestamps
}

func (ReportEntry) TableName() string { return "report_entries" }
