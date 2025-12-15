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

type projectStore struct {
	projects     map[uuid.UUID]*models.Project
	milestones   map[uuid.UUID]*models.Milestone
	tasks        map[uuid.UUID]*models.Task
	projectLinks map[uuid.UUID]*models.ProjectLink
	taskLinks    map[uuid.UUID]*models.TaskLink
}

// ProjectService keeps lightweight in-memory data to unblock UI flows and CRUD demos.
type ProjectService struct {
	mu     sync.RWMutex
	stores map[uuid.UUID]*projectStore
}

func NewProjectService() *ProjectService {
	return &ProjectService{stores: make(map[uuid.UUID]*projectStore)}
}

func (s *ProjectService) ensureStore(orgID uuid.UUID) *projectStore {
	s.mu.Lock()
	defer s.mu.Unlock()
	store, ok := s.stores[orgID]
	if ok {
		return store
	}
	store = &projectStore{
		projects:     make(map[uuid.UUID]*models.Project),
		milestones:   make(map[uuid.UUID]*models.Milestone),
		tasks:        make(map[uuid.UUID]*models.Task),
		projectLinks: make(map[uuid.UUID]*models.ProjectLink),
		taskLinks:    make(map[uuid.UUID]*models.TaskLink),
	}

	// seed demo data
	projectID := uuid.New()
	now := time.Now()
	store.projects[projectID] = &models.Project{
		ID:            projectID,
		OrgID:         orgID,
		OwnerUserID:   uuid.New(),
		Name:          "Unified Workspace Migration",
		Description:   "Move project execution flows into WorkPulse with OKR alignment.",
		Status:        "active",
		StartDate:     timePtr(now.AddDate(0, -1, 0)),
		EndDate:       timePtr(now.AddDate(0, 1, 0)),
		SchemaVersion: 1,
		Payload:       datatypes.JSON([]byte(`{"theme":"execution"}`)),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	ms1 := uuid.New()
	store.milestones[ms1] = &models.Milestone{
		ID:          ms1,
		OrgID:       orgID,
		ProjectID:   projectID,
		Title:       "Backend foundation",
		Description: "CRUD APIs for projects, milestones, and tasks",
		DueDate:     timePtr(now.AddDate(0, 0, 21)),
		Status:      "active",
		Progress:    0,
	}

	ms2 := uuid.New()
	store.milestones[ms2] = &models.Milestone{
		ID:          ms2,
		OrgID:       orgID,
		ProjectID:   projectID,
		Title:       "Collaboration preview",
		Description: "Ship Gantt and Kanban views",
		DueDate:     timePtr(now.AddDate(0, 1, 5)),
		Status:      "active",
		Progress:    0,
	}

	task1 := uuid.New()
	store.tasks[task1] = &models.Task{
		ID:          task1,
		OrgID:       orgID,
		ProjectID:   projectID,
		MilestoneID: &ms1,
		Title:       "Model project entities",
		Status:      "doing",
		Tags:        datatypes.JSON([]byte(`["backend"]`)),
		Order:       1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Metadata:    datatypes.JSON([]byte(`{}`)),
	}

	task2 := uuid.New()
	store.tasks[task2] = &models.Task{
		ID:           task2,
		OrgID:        orgID,
		ProjectID:    projectID,
		MilestoneID:  &ms1,
		ParentTaskID: &task1,
		Title:        "Wire Gin handlers",
		Status:       "todo",
		Tags:         datatypes.JSON([]byte(`["backend","api"]`)),
		Order:        2,
		CreatedAt:    now,
		UpdatedAt:    now,
		Metadata:     datatypes.JSON([]byte(`{}`)),
	}

	task3 := uuid.New()
	store.tasks[task3] = &models.Task{
		ID:          task3,
		OrgID:       orgID,
		ProjectID:   projectID,
		MilestoneID: &ms2,
		Title:       "Render Kanban board",
		Status:      "doing",
		Tags:        datatypes.JSON([]byte(`["frontend","ui"]`)),
		Order:       1,
		CreatedAt:   now,
		UpdatedAt:   now,
		Metadata:    datatypes.JSON([]byte(`{}`)),
	}

	linkID := uuid.New()
	store.projectLinks[linkID] = &models.ProjectLink{
		ID:         linkID,
		OrgID:      orgID,
		ProjectID:  projectID,
		EntityType: "okr_kr",
		EntityID:   uuid.New(),
		Relation:   "supports",
		CreatedAt:  now,
	}

	s.stores[orgID] = store
	return store
}

func timePtr(t time.Time) *time.Time { return &t }

// ProjectWithStats wraps derived counts and progress.
type ProjectWithStats struct {
	*models.Project
	MilestoneCount int     `json:"milestone_count"`
	TaskCount      int     `json:"task_count"`
	Progress       float64 `json:"progress"`
}

// ListProjects returns projects with derived progress and counts.
func (s *ProjectService) ListProjects(ctx context.Context, orgID uuid.UUID) ([]ProjectWithStats, error) {
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []ProjectWithStats
	for _, p := range store.projects {
		milestones := s.milestonesForProjectLocked(store, p.ID)
		tasks := s.tasksForProjectLocked(store, p.ID)
		progress := s.computeProjectProgressLocked(milestones, tasks)
		res = append(res, ProjectWithStats{
			Project:        p,
			MilestoneCount: len(milestones),
			TaskCount:      len(tasks),
			Progress:       progress,
		})
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].UpdatedAt.After(res[j].UpdatedAt)
	})
	return res, nil
}

func (s *ProjectService) milestonesForProjectLocked(store *projectStore, projectID uuid.UUID) []*models.Milestone {
	var res []*models.Milestone
	for _, m := range store.milestones {
		if m.ProjectID == projectID {
			clone := *m
			res = append(res, &clone)
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i].Title < res[j].Title })
	return res
}

func (s *ProjectService) tasksForProjectLocked(store *projectStore, projectID uuid.UUID) []*models.Task {
	var res []*models.Task
	for _, t := range store.tasks {
		if t.ProjectID == projectID {
			clone := *t
			res = append(res, &clone)
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i].Order < res[j].Order })
	return res
}

func (s *ProjectService) computeProjectProgressLocked(milestones []*models.Milestone, tasks []*models.Task) float64 {
	if len(milestones) == 0 && len(tasks) == 0 {
		return 0
	}
	var total float64
	var count int
	milestoneProgress := s.computeMilestoneProgressLocked(milestones, tasks)
	for _, p := range milestoneProgress {
		total += p
		count++
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func (s *ProjectService) computeMilestoneProgressLocked(milestones []*models.Milestone, tasks []*models.Task) map[uuid.UUID]float64 {
	result := make(map[uuid.UUID]float64)
	tasksByMilestone := make(map[uuid.UUID][]*models.Task)
	for _, t := range tasks {
		if t.MilestoneID != nil {
			tasksByMilestone[*t.MilestoneID] = append(tasksByMilestone[*t.MilestoneID], t)
		}
	}

	for _, m := range milestones {
		lst := tasksByMilestone[m.ID]
		if len(lst) == 0 {
			result[m.ID] = 0
			continue
		}
		var done int
		for _, t := range lst {
			if t.Status == "done" {
				done++
			}
		}
		result[m.ID] = float64(done) / float64(len(lst)) * 100
	}
	return result
}

// CreateProject persists a new project.
func (s *ProjectService) CreateProject(ctx context.Context, p *models.Project) error {
	store := s.ensureStore(p.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.projects[p.ID] = p
	return nil
}

// UpdateProject updates an existing project in place.
func (s *ProjectService) UpdateProject(ctx context.Context, orgID, id uuid.UUID, name, desc, status string, start, end *time.Time) (*models.Project, error) {
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := store.projects[id]
	if !ok {
		return nil, nil
	}
	if name != "" {
		p.Name = name
	}
	p.Description = desc
	if status != "" {
		p.Status = status
	}
	p.StartDate = start
	p.EndDate = end
	p.UpdatedAt = time.Now()
	return p, nil
}

// CreateMilestone creates a milestone under a project.
func (s *ProjectService) CreateMilestone(ctx context.Context, m *models.Milestone) error {
	store := s.ensureStore(m.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.milestones[m.ID] = m
	return nil
}

// UpdateMilestone updates milestone details.
func (s *ProjectService) UpdateMilestone(ctx context.Context, orgID, id uuid.UUID, title, desc, status string, due *time.Time) (*models.Milestone, error) {
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := store.milestones[id]
	if !ok {
		return nil, nil
	}
	if title != "" {
		m.Title = title
	}
	m.Description = desc
	if status != "" {
		m.Status = status
	}
	m.DueDate = due
	return m, nil
}

// ListMilestones returns milestones with progress.
func (s *ProjectService) ListMilestones(ctx context.Context, orgID, projectID uuid.UUID) ([]*models.Milestone, error) {
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	milestones := s.milestonesForProjectLocked(store, projectID)
	tasks := s.tasksForProjectLocked(store, projectID)
	progressMap := s.computeMilestoneProgressLocked(milestones, tasks)
	for _, m := range milestones {
		m.Progress = progressMap[m.ID]
	}
	return milestones, nil
}

// CreateTask creates a task in a project.
func (s *ProjectService) CreateTask(ctx context.Context, t *models.Task) error {
	store := s.ensureStore(t.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.tasks[t.ID] = t
	return nil
}

// UpdateTask patches a task.
func (s *ProjectService) UpdateTask(ctx context.Context, orgID, id uuid.UUID, title, status string, tags datatypes.JSON, milestoneID, parentID *uuid.UUID) (*models.Task, error) {
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := store.tasks[id]
	if !ok {
		return nil, nil
	}
	if title != "" {
		task.Title = title
	}
	if status != "" {
		task.Status = status
	}
	if tags != nil {
		task.Tags = tags
	}
	if milestoneID != nil {
		task.MilestoneID = milestoneID
	}
	if parentID != nil {
		task.ParentTaskID = parentID
	}
	task.UpdatedAt = time.Now()
	return task, nil
}

// ListTasks returns tasks filtered by project.
func (s *ProjectService) ListTasks(ctx context.Context, orgID, projectID uuid.UUID) ([]*models.Task, error) {
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := s.tasksForProjectLocked(store, projectID)
	return tasks, nil
}

// AddProjectLink links a project to an external artifact.
func (s *ProjectService) AddProjectLink(ctx context.Context, link *models.ProjectLink) error {
	store := s.ensureStore(link.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.projectLinks[link.ID] = link
	return nil
}

// RemoveProjectLink deletes a project link.
func (s *ProjectService) RemoveProjectLink(ctx context.Context, orgID, id uuid.UUID) {
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(store.projectLinks, id)
}

// AddTaskLink adds a link from a task to OKR or report.
func (s *ProjectService) AddTaskLink(ctx context.Context, link *models.TaskLink) error {
	store := s.ensureStore(link.OrgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	store.taskLinks[link.ID] = link
	return nil
}

// RemoveTaskLink deletes a task link.
func (s *ProjectService) RemoveTaskLink(ctx context.Context, orgID, id uuid.UUID) {
	store := s.ensureStore(orgID)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(store.taskLinks, id)
}

// ListTaskLinks returns linked artifacts for a task.
func (s *ProjectService) ListTaskLinks(ctx context.Context, orgID, taskID uuid.UUID) ([]*models.TaskLink, error) {
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	var res []*models.TaskLink
	for _, l := range store.taskLinks {
		if l.TaskID == taskID {
			clone := *l
			res = append(res, &clone)
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i].CreatedAt.After(res[j].CreatedAt) })
	return res, nil
}

// ProjectOKRProgress aggregates milestone progress into linked key results.
type ProjectOKRProgress struct {
	ProjectProgress float64 `json:"project_progress"`
	KeyResults      []struct {
		KeyResultID uuid.UUID `json:"key_result_id"`
		Progress    float64   `json:"progress"`
	} `json:"key_results"`
}

// ComputeOKRProgress aggregates completed milestone ratio into KR progress.
func (s *ProjectService) ComputeOKRProgress(ctx context.Context, orgID, projectID uuid.UUID) (ProjectOKRProgress, error) {
	store := s.ensureStore(orgID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	milestones := s.milestonesForProjectLocked(store, projectID)
	tasks := s.tasksForProjectLocked(store, projectID)
	milestoneProgress := s.computeMilestoneProgressLocked(milestones, tasks)
	overall := s.computeProjectProgressLocked(milestones, tasks)

	var result ProjectOKRProgress
	result.ProjectProgress = overall
	for _, l := range store.projectLinks {
		if l.ProjectID != projectID || l.EntityType != "okr_kr" {
			continue
		}
		progress := float64(0)
		if len(milestoneProgress) > 0 {
			for _, p := range milestoneProgress {
				progress += p
			}
			progress = progress / float64(len(milestoneProgress))
		}
		result.KeyResults = append(result.KeyResults, struct {
			KeyResultID uuid.UUID `json:"key_result_id"`
			Progress    float64   `json:"progress"`
		}{KeyResultID: l.EntityID, Progress: progress})
	}
	return result, nil
}
