package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Project is a light-weight aggregate for delivery initiatives.
type Project struct {
	ID            uuid.UUID      `json:"id"`
	OrgID         uuid.UUID      `json:"org_id"`
	TeamID        *uuid.UUID     `json:"team_id"`
	OwnerUserID   uuid.UUID      `json:"owner_user_id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Status        string         `json:"status"`
	StartDate     *time.Time     `json:"start_date"`
	EndDate       *time.Time     `json:"end_date"`
	SchemaVersion int            `json:"schema_version"`
	Payload       datatypes.JSON `json:"payload"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// Milestone groups a slice of work items within a project.
type Milestone struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Status      string     `json:"status"`
	Progress    float64    `json:"progress"`
}

// Task is a project-scoped unit of work with optional parent/child relationship.
type Task struct {
	ID           uuid.UUID      `json:"id"`
	OrgID        uuid.UUID      `json:"org_id"`
	ProjectID    uuid.UUID      `json:"project_id"`
	MilestoneID  *uuid.UUID     `json:"milestone_id"`
	ParentTaskID *uuid.UUID     `json:"parent_task_id"`
	Title        string         `json:"title"`
	Status       string         `json:"status"`
	Tags         datatypes.JSON `json:"tags"`
	AssigneeID   *uuid.UUID     `json:"assignee_id"`
	Order        int            `json:"order"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Metadata     datatypes.JSON `json:"metadata"`
}

// ProjectLink connects a project to OKRs, reports, or meetings.
type ProjectLink struct {
	ID         uuid.UUID `json:"id"`
	OrgID      uuid.UUID `json:"org_id"`
	ProjectID  uuid.UUID `json:"project_id"`
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
	Relation   string    `json:"relation"`
	CreatedAt  time.Time `json:"created_at"`
}

// TaskLink connects a task to external artifacts like OKRs or reports.
type TaskLink struct {
	ID         uuid.UUID `json:"id"`
	OrgID      uuid.UUID `json:"org_id"`
	TaskID     uuid.UUID `json:"task_id"`
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
	CreatedAt  time.Time `json:"created_at"`
}
