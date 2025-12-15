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

// OKRCoverage captures how many key results are linked to a specific entity type.
type OKRCoverage struct {
	EntityType       string  `json:"entity_type"`
	Coverage         float64 `json:"coverage"`
	LinkedKeyResults int     `json:"linked_key_results"`
	TotalKeyResults  int     `json:"total_key_results"`
	Deviation        float64 `json:"deviation"`
}

// OKRMetrics aggregates completion, deviation, and coverage insights.
type OKRMetrics struct {
	Completion         float64               `json:"completion"`
	Deviation          float64               `json:"deviation"`
	UnlinkedKeyResults []models.OKRKeyResult `json:"unlinked_key_results"`
	Coverage           []OKRCoverage         `json:"coverage"`
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

// ComputeMetrics aggregates completion, deviation and link coverage for the selected slice of OKRs.
func (s *OKRService) ComputeMetrics(ctx context.Context, orgID uuid.UUID, cycleID, teamID, ownerID *uuid.UUID) (OKRMetrics, error) {
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filteredObjectives []*models.OKRObjective
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
		filteredObjectives = append(filteredObjectives, obj)
	}

	objectiveLookup := make(map[uuid.UUID]struct{}, len(filteredObjectives))
	for _, obj := range filteredObjectives {
		objectiveLookup[obj.ID] = struct{}{}
	}

	var keyResults []*models.OKRKeyResult
	for _, kr := range store.keyResults {
		if _, ok := objectiveLookup[kr.ObjectiveID]; !ok {
			continue
		}
		keyResults = append(keyResults, kr)
	}

	metrics := OKRMetrics{}
	if len(keyResults) == 0 {
		return metrics, nil
	}

	var totalProgress float64
	var totalDeviation float64
	linkedByType := make(map[string]map[uuid.UUID]struct{})
	progressByKR := make(map[uuid.UUID]float64)

	for _, kr := range keyResults {
		progress := progressValue(kr)
		progressByKR[kr.ID] = progress
		totalProgress += progress
		totalDeviation += (100 - progress)

		var hasLink bool
		for _, l := range store.links {
			if l.KeyResultID != nil && *l.KeyResultID == kr.ID {
				hasLink = true
				if _, ok := linkedByType[l.EntityType]; !ok {
					linkedByType[l.EntityType] = make(map[uuid.UUID]struct{})
				}
				linkedByType[l.EntityType][kr.ID] = struct{}{}
			}
		}
		if !hasLink {
			metrics.UnlinkedKeyResults = append(metrics.UnlinkedKeyResults, *kr)
		}
	}

	metrics.Completion = totalProgress / float64(len(keyResults))
	metrics.Deviation = totalDeviation / float64(len(keyResults))

	for entityType, linked := range linkedByType {
		linkedCount := len(linked)
		coverage := (float64(linkedCount) / float64(len(keyResults))) * 100
		var devSum float64
		for krID := range linked {
			devSum += (100 - progressByKR[krID])
		}
		deviation := 0.0
		if linkedCount > 0 {
			deviation = devSum / float64(linkedCount)
		}
		metrics.Coverage = append(metrics.Coverage, OKRCoverage{
			EntityType:       entityType,
			Coverage:         coverage,
			LinkedKeyResults: linkedCount,
			TotalKeyResults:  len(keyResults),
			Deviation:        deviation,
		})
	}

	return metrics, nil
}

func progressValue(kr *models.OKRKeyResult) float64 {
	if kr.TargetValue == nil || *kr.TargetValue == 0 || kr.CurrentValue == nil {
		return 0
	}
	progress := (*kr.CurrentValue / *kr.TargetValue) * 100
	if progress > 100 {
		return 100
	}
	return progress
}
