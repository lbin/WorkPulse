package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"workpulse/internal/models"
	"workpulse/internal/repo"
)

type ReportService struct {
	repo      *repo.ReportRepo
	auditRepo *repo.AuditRepo
}

func NewReportService(repo *repo.ReportRepo, auditRepo *repo.AuditRepo) *ReportService {
	return &ReportService{repo: repo, auditRepo: auditRepo}
}

func (s *ReportService) Create(ctx context.Context, actor uuid.UUID, report *models.Report, links []models.ReportLink) error {
	if report.SchemaVersion == 0 {
		report.SchemaVersion = 1
	}
	if len(report.Content) == 0 {
		report.Content = datatypes.JSON([]byte(`{}`))
	}
	if len(report.Payload) == 0 {
		report.Payload = datatypes.JSON([]byte(`{}`))
	}
	if report.Status == "" {
		report.Status = "draft"
	}
	if err := s.repo.Create(ctx, report, links); err != nil {
		return err
	}
	return s.auditRepo.Add(ctx, models.AuditLog{
		OrgID:       report.OrgID,
		ActorUserID: actor,
		Action:      "report.create",
		EntityType:  "report",
		EntityID:    report.ID,
	})
}

func (s *ReportService) Update(ctx context.Context, actor uuid.UUID, report *models.Report, links []models.ReportLink) error {
	current, _, err := s.repo.Get(ctx, report.OrgID, report.ID)
	if err != nil {
		return err
	}
	if current.Status == "approved" || current.Status == "archived" {
		return errors.New("approved/archived reports cannot be edited")
	}
	if report.SchemaVersion == 0 {
		report.SchemaVersion = current.SchemaVersion
	}
	if len(report.Content) == 0 {
		report.Content = current.Content
	}
	if err := s.repo.Update(ctx, report, links); err != nil {
		return err
	}
	return s.auditRepo.Add(ctx, models.AuditLog{
		OrgID:       report.OrgID,
		ActorUserID: actor,
		Action:      "report.update",
		EntityType:  "report",
		EntityID:    report.ID,
	})
}

func (s *ReportService) Get(ctx context.Context, orgID, id uuid.UUID) (*models.Report, []models.ReportLink, error) {
	return s.repo.Get(ctx, orgID, id)
}

func (s *ReportService) List(ctx context.Context, orgID uuid.UUID, filter repo.ReportFilter) ([]models.Report, int64, error) {
	return s.repo.List(ctx, orgID, filter)
}

func (s *ReportService) ChangeStatus(ctx context.Context, orgID, reportID, actor uuid.UUID, action string, comment *string) (*models.Report, error) {
	report, _, err := s.repo.Get(ctx, orgID, reportID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	switch action {
	case "submit":
		if report.Status != "draft" && report.Status != "rejected" {
			return nil, errors.New("only draft/rejected reports can be submitted")
		}
		if err := s.repo.UpdateStatus(ctx, orgID, reportID, "submitted", nil, nil, &now, nil); err != nil {
			return nil, err
		}
		report.Status = "submitted"
		report.SubmittedAt = &now
	case "approve":
		if report.Status != "submitted" {
			return nil, errors.New("only submitted reports can be approved")
		}
		if err := s.repo.UpdateStatus(ctx, orgID, reportID, "approved", &actor, comment, report.SubmittedAt, &now); err != nil {
			return nil, err
		}
		report.Status = "approved"
		report.ReviewerUserID = &actor
		report.ReviewComment = comment
		report.ReviewedAt = &now
	case "reject":
		if report.Status != "submitted" {
			return nil, errors.New("only submitted reports can be rejected")
		}
		if err := s.repo.UpdateStatus(ctx, orgID, reportID, "rejected", &actor, comment, report.SubmittedAt, &now); err != nil {
			return nil, err
		}
		report.Status = "rejected"
		report.ReviewerUserID = &actor
		report.ReviewComment = comment
		report.ReviewedAt = &now
	default:
		return nil, errors.New("unknown action")
	}

	diffBytes, _ := json.Marshal(map[string]string{"status": report.Status})
	_ = s.auditRepo.Add(ctx, models.AuditLog{
		OrgID:       orgID,
		ActorUserID: actor,
		Action:      "report." + action,
		EntityType:  "report",
		EntityID:    reportID,
		Diff:        datatypes.JSON(diffBytes),
	})

	return report, nil
}
