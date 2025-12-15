package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"workpulse/internal/models"
)

type ReportFilter struct {
	Type        *string
	Status      *string
	AuthorID    *uuid.UUID
	PeriodStart *time.Time
	PeriodEnd   *time.Time
	Limit       int
	Offset      int
}

type ReportRepo struct{ db *gorm.DB }

func NewReportRepo(db *gorm.DB) *ReportRepo { return &ReportRepo{db: db} }

func (r *ReportRepo) Create(ctx context.Context, report *models.Report, links []models.ReportLink) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(report).Error; err != nil {
			return err
		}
		if len(links) > 0 {
			for i := range links {
				links[i].OrgID = report.OrgID
				links[i].ReportID = report.ID
			}
			if err := tx.Create(&links).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ReportRepo) Update(ctx context.Context, report *models.Report, links []models.ReportLink) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Report{}).Where("id=? AND org_id=?", report.ID, report.OrgID).
			Updates(map[string]interface{}{
				"type":         report.Type,
				"period_start": report.PeriodStart,
				"period_end":   report.PeriodEnd,
				"title":        report.Title,
				"summary":      report.Summary,
				"content":      report.Content,
				"payload":      report.Payload,
				"updated_at":   time.Now(),
			}).Error; err != nil {
			return err
		}
		if err := r.replaceLinks(ctx, tx, report.OrgID, report.ID, links); err != nil {
			return err
		}
		return nil
	})
}

func (r *ReportRepo) replaceLinks(ctx context.Context, tx *gorm.DB, orgID, reportID uuid.UUID, links []models.ReportLink) error {
	if err := tx.WithContext(ctx).Where("org_id=? AND report_id=?", orgID, reportID).Delete(&models.ReportLink{}).Error; err != nil {
		return err
	}
	if len(links) == 0 {
		return nil
	}
	for i := range links {
		links[i].OrgID = orgID
		links[i].ReportID = reportID
	}
	return tx.Create(&links).Error
}

func (r *ReportRepo) UpdateStatus(ctx context.Context, orgID, reportID uuid.UUID, status string, reviewer *uuid.UUID, comment *string, submittedAt *time.Time, reviewedAt *time.Time) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if reviewer != nil {
		updates["reviewer_user_id"] = reviewer
	}
	if comment != nil {
		updates["review_comment"] = comment
	}
	if submittedAt != nil {
		updates["submitted_at"] = submittedAt
	}
	if reviewedAt != nil {
		updates["reviewed_at"] = reviewedAt
	}
	updates["updated_at"] = time.Now()

	return r.db.WithContext(ctx).Model(&models.Report{}).
		Where("org_id=? AND id=?", orgID, reportID).
		Updates(updates).Error
}

func (r *ReportRepo) Get(ctx context.Context, orgID, id uuid.UUID) (*models.Report, []models.ReportLink, error) {
	var rep models.Report
	if err := r.db.WithContext(ctx).Where("org_id=? AND id=? AND deleted_at IS NULL", orgID, id).First(&rep).Error; err != nil {
		return nil, nil, err
	}
	var links []models.ReportLink
	if err := r.db.WithContext(ctx).Where("org_id=? AND report_id=?", orgID, id).Order("created_at").Find(&links).Error; err != nil {
		return nil, nil, err
	}
	return &rep, links, nil
}

func (r *ReportRepo) List(ctx context.Context, orgID uuid.UUID, filter ReportFilter) ([]models.Report, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Report{}).Where("org_id=? AND deleted_at IS NULL", orgID)
	if filter.Type != nil {
		q = q.Where("type=?", *filter.Type)
	}
	if filter.Status != nil {
		q = q.Where("status=?", *filter.Status)
	}
	if filter.AuthorID != nil {
		q = q.Where("author_user_id=?", *filter.AuthorID)
	}
	if filter.PeriodStart != nil {
		q = q.Where("period_start >= ?", *filter.PeriodStart)
	}
	if filter.PeriodEnd != nil {
		q = q.Where("period_end <= ?", *filter.PeriodEnd)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	limit := filter.Limit
	if limit == 0 {
		limit = 20
	}
	var reports []models.Report
	if err := q.Order("period_start DESC, updated_at DESC").Limit(limit).Offset(filter.Offset).Find(&reports).Error; err != nil {
		return nil, 0, err
	}
	return reports, total, nil
}
