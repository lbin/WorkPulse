package service

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"workpulse/internal/models"
)

type meetingStore struct {
	meetings map[uuid.UUID]*models.Meeting
	actions  map[uuid.UUID]*models.MeetingAction
	links    map[uuid.UUID]*models.MeetingLink
}

// MeetingService maintains simple in-memory meetings data to support UI flows.
type MeetingService struct {
	mu     sync.RWMutex
	stores map[uuid.UUID]*meetingStore
}

func NewMeetingService() *MeetingService {
	return &MeetingService{stores: make(map[uuid.UUID]*meetingStore)}
}

func (s *MeetingService) ensureStore(orgID uuid.UUID) *meetingStore {
	s.mu.Lock()
	defer s.mu.Unlock()
	store, ok := s.stores[orgID]
	if ok {
		return store
	}
	now := time.Now()
	meetingID := uuid.New()
	store = &meetingStore{
		meetings: map[uuid.UUID]*models.Meeting{
			meetingID: {
				ID:                meetingID,
				OrgID:             orgID,
				Title:             "Weekly Execution Sync",
				Agenda:            "Review OKR progress, unblock tasks, plan next sprint items",
				ScheduledAt:       now.AddDate(0, 0, 1),
				DurationMinutes:   45,
				FacilitatorUserID: nil,
				Notes:             "Discussed project risk on integration timeline.",
				AttendeeIDs:       datatypes.JSON([]byte(`["u1","u2","u3"]`)),
				Status:            "scheduled",
				SchemaVersion:     1,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		actions: make(map[uuid.UUID]*models.MeetingAction),
		links:   make(map[uuid.UUID]*models.MeetingLink),
	}
	actionID := uuid.New()
	store.actions[actionID] = &models.MeetingAction{
		ID:        actionID,
		OrgID:     orgID,
		MeetingID: meetingID,
		Title:     "Validate data contract with analytics team",
		Status:    "todo",
		CreatedAt: now,
		UpdatedAt: now,
	}
	linkID := uuid.New()
	store.links[linkID] = &models.MeetingLink{
		ID:         linkID,
		OrgID:      orgID,
		MeetingID:  meetingID,
		TargetType: "project",
		TargetID:   uuid.New(),
		Relation:   "related",
		CreatedAt:  now,
	}

	s.stores[orgID] = store
	return store
}

func (s *MeetingService) List(ctx context.Context, orgID uuid.UUID, teamID *uuid.UUID, start *time.Time, end *time.Time) ([]*models.Meeting, error) {
	_ = ctx
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*models.Meeting
	for _, m := range store.meetings {
		if teamID != nil && (m.TeamID == nil || *m.TeamID != *teamID) {
			continue
		}
		if start != nil && m.ScheduledAt.Before(*start) {
			continue
		}
		if end != nil && m.ScheduledAt.After(*end) {
			continue
		}
		copy := *m
		result = append(result, &copy)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ScheduledAt.After(result[j].ScheduledAt) })
	return result, nil
}

func (s *MeetingService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Meeting, []*models.MeetingAction, []*models.MeetingLink, error) {
	_ = ctx
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	meeting, ok := store.meetings[id]
	if !ok {
		return nil, nil, nil, nil
	}
	mCopy := *meeting
	var actions []*models.MeetingAction
	for _, a := range store.actions {
		if a.MeetingID == id {
			c := *a
			actions = append(actions, &c)
		}
	}
	sort.Slice(actions, func(i, j int) bool { return actions[i].CreatedAt.After(actions[j].CreatedAt) })
	var links []*models.MeetingLink
	for _, l := range store.links {
		if l.MeetingID == id {
			c := *l
			links = append(links, &c)
		}
	}
	return &mCopy, actions, links, nil
}

func (s *MeetingService) Create(ctx context.Context, meeting *models.Meeting) error {
	_ = ctx
	store := s.ensureStore(meeting.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.meetings[meeting.ID] = meeting
	return nil
}

func (s *MeetingService) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, update func(m *models.Meeting)) (*models.Meeting, error) {
	_ = ctx
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	meeting, ok := store.meetings[id]
	if !ok {
		return nil, nil
	}
	update(meeting)
	meeting.UpdatedAt = time.Now()
	copy := *meeting
	return &copy, nil
}

func (s *MeetingService) AddAction(ctx context.Context, action *models.MeetingAction) error {
	_ = ctx
	store := s.ensureStore(action.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.actions[action.ID] = action
	return nil
}

func (s *MeetingService) UpdateAction(ctx context.Context, orgID, actionID uuid.UUID, fn func(a *models.MeetingAction)) (*models.MeetingAction, error) {
	_ = ctx
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	action, ok := store.actions[actionID]
	if !ok {
		return nil, nil
	}
	fn(action)
	action.UpdatedAt = time.Now()
	copy := *action
	return &copy, nil
}

func (s *MeetingService) AddLink(ctx context.Context, link *models.MeetingLink) error {
	_ = ctx
	store := s.ensureStore(link.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.links[link.ID] = link
	return nil
}
