package models

import "github.com/google/uuid"

type Org struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name   string
	Status string
	Timestamps
	SoftDelete
}

func (Org) TableName() string { return "orgs" }
