package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Report struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID          uuid.UUID  `gorm:"type:uuid;index"`
	TeamID         *uuid.UUID `gorm:"type:uuid;index"`
	AuthorUserID   uuid.UUID  `gorm:"type:uuid;index;column:author_user_id"`
	Type           string
	PeriodStart    time.Time `gorm:"type:date;column:period_start"`
	PeriodEnd      time.Time `gorm:"type:date;column:period_end"`
	Status         string
	Title          *string
	Summary        *string
	SchemaVersion  int            `gorm:"column:schema_version"`
	Content        datatypes.JSON `gorm:"type:jsonb"`
	Payload        datatypes.JSON `gorm:"type:jsonb"`
	SubmittedAt    *time.Time     `gorm:"column:submitted_at"`
	ReviewedAt     *time.Time     `gorm:"column:reviewed_at"`
	ReviewerUserID *uuid.UUID     `gorm:"type:uuid;column:reviewer_user_id"`
	ReviewComment  *string        `gorm:"column:review_comment"`
	Timestamps
	SoftDelete
}

func (Report) TableName() string { return "reports" }

type ReportEntry struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID         uuid.UUID `gorm:"type:uuid;index"`
	ReportID      uuid.UUID `gorm:"type:uuid;index;column:report_id"`
	Section       string
	Content       string
	EffortMinutes *int `gorm:"column:effort_minutes"`
	Status        string
	WorkItemID    *uuid.UUID `gorm:"type:uuid;index;column:work_item_id"`
	OrderNo       int        `gorm:"column:order_no"`
	Timestamps
}

func (ReportEntry) TableName() string { return "report_entries" }

type ReportLink struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID      uuid.UUID `gorm:"type:uuid;index"`
	ReportID   uuid.UUID `gorm:"type:uuid;index;column:report_id"`
	TargetType string
	TargetID   uuid.UUID `gorm:"type:uuid;index;column:target_id"`
	Relation   string
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ReportLink) TableName() string { return "report_links" }
