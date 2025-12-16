package repo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"workpulse/internal/models"
)

// UserRepo handles persistence for user entities.
type UserRepo struct{ db *gorm.DB }

// NewUserRepo builds a user repository backed by gorm.
func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

// Create inserts a new user record.
func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// FindByEmail searches for a user by email within an organization.
func (r *UserRepo) FindByEmail(ctx context.Context, orgID uuid.UUID, email string) (*models.User, error) {
	var u models.User
	q := r.db.WithContext(ctx).Where("email=? AND deleted_at IS NULL", email)
	if orgID != uuid.Nil {
		q = q.Where("org_id=?", orgID)
	}
	err := q.First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Get fetches a user by ID scoped to organization.
func (r *UserRepo) Get(ctx context.Context, orgID, userID uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("org_id=? AND id=? AND deleted_at IS NULL", orgID, userID).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
