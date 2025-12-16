package repo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"workpulse/internal/models"
)

// MembershipRepo manages relationships between users and teams.
type MembershipRepo struct{ db *gorm.DB }

// NewMembershipRepo constructs a repository instance.
func NewMembershipRepo(db *gorm.DB) *MembershipRepo { return &MembershipRepo{db: db} }

// Create adds a membership row for the given team and user.
func (r *MembershipRepo) Create(ctx context.Context, m *models.Membership) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ListByUser returns all memberships for a specific user.
func (r *MembershipRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Membership, error) {
	var ms []models.Membership
	err := r.db.WithContext(ctx).Where("user_id=? AND deleted_at IS NULL", userID).Find(&ms).Error
	return ms, err
}

// Exists checks if the given user already belongs to the target team.
func (r *MembershipRepo) Exists(ctx context.Context, teamID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Membership{}).Where("team_id=? AND user_id=? AND deleted_at IS NULL", teamID, userID).Count(&count).Error
	return count > 0, err
}
