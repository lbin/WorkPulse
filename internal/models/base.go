package models

import "time"

type SoftDelete struct {
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

type Timestamps struct {
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
