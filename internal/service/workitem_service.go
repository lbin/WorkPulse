package service

import (
	"context"

	"github.com/google/uuid"
	"workpulse/internal/models"
	"workpulse/internal/repo"
)

type WorkItemService struct {
	wiRepo *repo.WorkItemRepo
}

func NewWorkItemService(wiRepo *repo.WorkItemRepo) *WorkItemService {
	return &WorkItemService{wiRepo: wiRepo}
}

func (s *WorkItemService) Create(ctx context.Context, wi *models.WorkItem) error {
	return s.wiRepo.Create(ctx, wi)
}

func (s *WorkItemService) AddOKRLink(ctx context.Context, orgID uuid.UUID, workItemID uuid.UUID, objectiveID *uuid.UUID, krID *uuid.UUID, linkType string, w float64, evidence bool) error {
	link := &models.WorkItemOKRLink{
		OrgID:              orgID,
		WorkItemID:         workItemID,
		ObjectiveID:        objectiveID,
		KeyResultID:        krID,
		LinkType:           linkType,
		ContributionWeight: w,
		EvidenceRequired:   evidence,
	}
	return s.wiRepo.AddOKRLink(ctx, link)
}
