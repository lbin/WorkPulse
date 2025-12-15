package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Meeting captures a scheduled discussion with agenda, notes, and attendees.
type Meeting struct {
	ID                uuid.UUID      `json:"id"`
	OrgID             uuid.UUID      `json:"org_id"`
	TeamID            *uuid.UUID     `json:"team_id"`
	Title             string         `json:"title"`
	Agenda            string         `json:"agenda"`
	ScheduledAt       time.Time      `json:"scheduled_at"`
	DurationMinutes   int            `json:"duration_minutes"`
	FacilitatorUserID *uuid.UUID     `json:"facilitator_user_id"`
	Notes             string         `json:"notes"`
	AttendeeIDs       datatypes.JSON `json:"attendee_ids"`
	Status            string         `json:"status"`
	SchemaVersion     int            `json:"schema_version"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// MeetingAction represents a follow-up item captured during a meeting.
type MeetingAction struct {
	ID            uuid.UUID  `json:"id"`
	OrgID         uuid.UUID  `json:"org_id"`
	MeetingID     uuid.UUID  `json:"meeting_id"`
	Title         string     `json:"title"`
	OwnerUserID   *uuid.UUID `json:"owner_user_id"`
	DueDate       *time.Time `json:"due_date"`
	Status        string     `json:"status"`
	RelatedTaskID *uuid.UUID `json:"related_task_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// MeetingLink connects a meeting to OKRs, projects, or reports.
type MeetingLink struct {
	ID         uuid.UUID `json:"id"`
	OrgID      uuid.UUID `json:"org_id"`
	MeetingID  uuid.UUID `json:"meeting_id"`
	TargetType string    `json:"target_type"`
	TargetID   uuid.UUID `json:"target_id"`
	Relation   string    `json:"relation"`
	CreatedAt  time.Time `json:"created_at"`
}
