package repo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"workpulse/internal/models"
)

type WorkItemRepo struct{ db *gorm.DB }

func NewWorkItemRepo(db *gorm.DB) *WorkItemRepo { return &WorkItemRepo{db: db} }

func (r *WorkItemRepo) Create(ctx context.Context, wi *models.WorkItem) error {
	return r.db.WithContext(ctx).Create(wi).Error
}

func (r *WorkItemRepo) Get(ctx context.Context, orgID, id uuid.UUID) (*models.WorkItem, error) {
	var wi models.WorkItem
	err := r.db.WithContext(ctx).Where("org_id=? AND id=? AND deleted_at IS NULL", orgID, id).First(&wi).Error
	if err != nil {
		return nil, err
	}
	return &wi, nil
}

func (r *WorkItemRepo) List(ctx context.Context, orgID uuid.UUID, teamID *uuid.UUID, ownerID *uuid.UUID, status *string, limit, offset int) ([]models.WorkItem, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.WorkItem{}).Where("org_id=? AND deleted_at IS NULL", orgID)
	if teamID != nil {
		q = q.Where("team_id=?", *teamID)
	}
	if ownerID != nil {
		q = q.Where("owner_user_id=?", *ownerID)
	}
	if status != nil {
		q = q.Where("status=?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.WorkItem
	err := q.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *WorkItemRepo) AddOKRLink(ctx context.Context, link *models.WorkItemOKRLink) error {
	return r.db.WithContext(ctx).Create(link).Error
}

func (r *WorkItemRepo) ListOKRLinks(ctx context.Context, orgID, workItemID uuid.UUID) ([]models.WorkItemOKRLink, error) {
	var links []models.WorkItemOKRLink
	err := r.db.WithContext(ctx).Where("org_id=? AND work_item_id=?", orgID, workItemID).Find(&links).Error
	return links, err
}
