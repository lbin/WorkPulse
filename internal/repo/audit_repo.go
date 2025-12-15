package repo

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"workpulse/internal/models"
)

type AuditRepo struct{ db *gorm.DB }

func NewAuditRepo(db *gorm.DB) *AuditRepo { return &AuditRepo{db: db} }

func (r *AuditRepo) Add(ctx context.Context, log models.AuditLog) error {
	if log.Diff == nil {
		log.Diff = datatypes.JSON([]byte(`{}`))
	}
	return r.db.WithContext(ctx).Create(&log).Error
}
