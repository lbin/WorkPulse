package repo

import (
	"context"

	"gorm.io/gorm"
	"workpulse/internal/models"
)

// OrgRepo provides persistence helpers for organizations.
type OrgRepo struct{ db *gorm.DB }

// NewOrgRepo constructs a repository for org entities.
func NewOrgRepo(db *gorm.DB) *OrgRepo { return &OrgRepo{db: db} }

// Create inserts a new organization record.
func (r *OrgRepo) Create(ctx context.Context, org *models.Org) error {
	return r.db.WithContext(ctx).Create(org).Error
}
