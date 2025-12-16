package repo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"workpulse/internal/models"
)

// TeamRepo wraps persistence for team entities.
type TeamRepo struct{ db *gorm.DB }

// NewTeamRepo builds a repository for teams.
func NewTeamRepo(db *gorm.DB) *TeamRepo { return &TeamRepo{db: db} }

// Create stores a new team in the database.
func (r *TeamRepo) Create(ctx context.Context, t *models.Team) error {
	return r.db.WithContext(ctx).Create(t).Error
}

// Get fetches a team by ID scoped to organization.
func (r *TeamRepo) Get(ctx context.Context, orgID, teamID uuid.UUID) (*models.Team, error) {
	var t models.Team
	err := r.db.WithContext(ctx).Where("org_id=? AND id=? AND deleted_at IS NULL", orgID, teamID).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// List returns all teams within an organization.
func (r *TeamRepo) List(ctx context.Context, orgID uuid.UUID) ([]models.Team, error) {
	var teams []models.Team
	err := r.db.WithContext(ctx).Where("org_id=? AND deleted_at IS NULL", orgID).Order("created_at ASC").Find(&teams).Error
	return teams, err
}
