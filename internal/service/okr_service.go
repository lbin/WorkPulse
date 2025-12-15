package service

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"workpulse/internal/models"
)

type okrStore struct {
	cycles     map[uuid.UUID]*models.OKRCycle
	objectives map[uuid.UUID]*models.OKRObjective
	keyResults map[uuid.UUID]*models.OKRKeyResult
	links      map[uuid.UUID]*models.OKRLink
}

// OKRService keeps a light-weight in-memory store so that the API can be exercised
// without requiring a full persistence layer. The storage is keyed per-org.
type OKRService struct {
	mu     sync.RWMutex
	stores map[uuid.UUID]*okrStore
}

func NewOKRService() *OKRService {
	return &OKRService{stores: make(map[uuid.UUID]*okrStore)}
}

func (s *OKRService) ensureStore(orgID uuid.UUID) *okrStore {
	s.mu.Lock()
	defer s.mu.Unlock()
	store, ok := s.stores[orgID]
	if ok {
		return store
	}
	store = &okrStore{
		cycles:     make(map[uuid.UUID]*models.OKRCycle),
		objectives: make(map[uuid.UUID]*models.OKRObjective),
		keyResults: make(map[uuid.UUID]*models.OKRKeyResult),
		links:      make(map[uuid.UUID]*models.OKRLink),
	}
	// seed with a default cycle and sample objective/KR so UI has data
	cycleID := uuid.New()
	store.cycles[cycleID] = &models.OKRCycle{
		ID:            cycleID,
		OrgID:         orgID,
		Type:          "quarter",
		Name:          "Q1 Preview",
		Status:        "active",
		SchemaVersion: 1,
	}

	objectiveID := uuid.New()
	store.objectives[objectiveID] = &models.OKRObjective{
		ID:            objectiveID,
		OrgID:         orgID,
		CycleID:       cycleID,
		Title:         "Ship collaboration launch",
		Status:        "active",
		SchemaVersion: 1,
		Tags:          datatypes.JSON([]byte(`[]`)),
		Payload:       datatypes.JSON([]byte(`{}`)),
	}

	krID := uuid.New()
	store.keyResults[krID] = &models.OKRKeyResult{
		ID:            krID,
		OrgID:         orgID,
		ObjectiveID:   objectiveID,
		Title:         "Reach 50 weekly active workspaces",
		MetricType:    "number",
		TargetValue:   floatPointer(50),
		CurrentValue:  floatPointer(10),
		Status:        "active",
		SchemaVersion: 1,
		MetricPayload: datatypes.JSON([]byte(`{}`)),
	}

	s.stores[orgID] = store
	return store
}

func floatPointer(v float64) *float64 { return &v }

func (s *OKRService) ListCycles(ctx context.Context, orgID uuid.UUID) ([]*models.OKRCycle, error) {
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*models.OKRCycle, 0, len(store.cycles))
	for _, c := range store.cycles {
		res = append(res, c)
	}
	return res, nil
}

func (s *OKRService) ListObjectives(ctx context.Context, orgID uuid.UUID, cycleID, teamID, ownerID *uuid.UUID, status *string) ([]models.OKRObjective, []models.OKRKeyResult, error) {
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	var objectives []models.OKRObjective
	var keyResults []models.OKRKeyResult
	for _, obj := range store.objectives {
		if cycleID != nil && obj.CycleID != *cycleID {
			continue
		}
		if teamID != nil && (obj.TeamID == nil || *obj.TeamID != *teamID) {
			continue
		}
		if ownerID != nil && obj.OwnerUserID != *ownerID {
			continue
		}
		if status != nil && obj.Status != *status {
			continue
		}
		objectives = append(objectives, *obj)
	}
	for _, kr := range store.keyResults {
		keyResults = append(keyResults, *kr)
	}
	return objectives, keyResults, nil
}

func (s *OKRService) CreateObjective(ctx context.Context, obj *models.OKRObjective) error {
	store := s.ensureStore(obj.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.objectives[obj.ID] = obj
	return nil
}

func (s *OKRService) UpdateObjective(ctx context.Context, obj *models.OKRObjective) error {
	store := s.ensureStore(obj.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.objectives[obj.ID] = obj
	return nil
}

func (s *OKRService) ArchiveObjective(ctx context.Context, orgID, id uuid.UUID) error {
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := store.objectives[id]; ok {
		existing.Status = "archived"
	}
	return nil
}

func (s *OKRService) CreateKeyResult(ctx context.Context, kr *models.OKRKeyResult) error {
	store := s.ensureStore(kr.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.keyResults[kr.ID] = kr
	return nil
}

func (s *OKRService) UpdateKeyResult(ctx context.Context, kr *models.OKRKeyResult) error {
	store := s.ensureStore(kr.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.keyResults[kr.ID] = kr
	return nil
}

func (s *OKRService) ArchiveKeyResult(ctx context.Context, orgID, id uuid.UUID) error {
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := store.keyResults[id]; ok {
		existing.Status = "archived"
	}
	return nil
}

func (s *OKRService) UpdateProgress(ctx context.Context, orgID, id uuid.UUID, current float64, confidence int) error {
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := store.keyResults[id]; ok {
		existing.CurrentValue = floatPointer(current)
		existing.Confidence = confidence
	}
	return nil
}

func (s *OKRService) AddLink(ctx context.Context, link *models.OKRLink) error {
	store := s.ensureStore(link.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.links[link.ID] = link
	return nil
}

func (s *OKRService) RemoveLink(ctx context.Context, orgID, id uuid.UUID) error {
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(store.links, id)
	return nil
}

func (s *OKRService) ListLinks(ctx context.Context, orgID uuid.UUID, entityType string, entityID uuid.UUID) ([]models.OKRLink, error) {
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	var res []models.OKRLink
	for _, l := range store.links {
		if l.EntityType == entityType && l.EntityID == entityID {
			res = append(res, *l)
		}
	}
	return res, nil
}
