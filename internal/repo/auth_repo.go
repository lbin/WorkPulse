package repo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"workpulse/internal/models"
)

// AuthRepo owns data-layer helpers for authentication and onboarding flows.
type AuthRepo struct{ db *gorm.DB }

// NewAuthRepo builds an AuthRepo using the provided database handle.
func NewAuthRepo(db *gorm.DB) *AuthRepo { return &AuthRepo{db: db} }

// FindUserByEmail returns a user record by email if it exists.
func (r *AuthRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUserWithOrg provisions a brand-new org and seeds its first user.
func (r *AuthRepo) CreateUserWithOrg(ctx context.Context, email, displayName, passwordHash string) (*models.User, *models.Org, error) {
	var user models.User
	var org models.Org

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		org = models.Org{Name: displayName, Status: "active"}
		if err := tx.Create(&org).Error; err != nil {
			return err
		}

		user = models.User{OrgID: org.ID, Email: email, DisplayName: displayName, Status: "active", PasswordHash: passwordHash}
		return tx.Create(&user).Error
	})

	if err != nil {
		return nil, nil, err
	}
	return &user, &org, nil
}

// CreateTeamWithMembership creates a team/org unit pair and subscribes the creator as owner.
func (r *AuthRepo) CreateTeamWithMembership(ctx context.Context, orgID, userID uuid.UUID, name string, parentTeamID *uuid.UUID) (*models.Team, error) {
	var team models.Team
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		team = models.Team{OrgID: orgID, ParentTeamID: parentTeamID, Name: name, Path: name}
		if err := tx.Create(&team).Error; err != nil {
			return err
		}

		orgUnit := models.OrgUnit{ID: team.ID, OrgID: orgID, ParentUnitID: parentTeamID, Name: name, Path: name, UnitType: "team"}
		if err := tx.Create(&orgUnit).Error; err != nil {
			return err
		}

		membership := models.Membership{OrgID: orgID, TeamID: &team.ID, OrgUnitID: &team.ID, UserID: userID, RoleInTeam: "owner"}
		return tx.Create(&membership).Error
	})

	if err != nil {
		return nil, err
	}
	return &team, nil
}

// AddUserToTeam subscribes a user to an existing team if not already present.
func (r *AuthRepo) AddUserToTeam(ctx context.Context, orgID, userID, teamID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var team models.Team
		if err := tx.Where("id=? AND org_id=? AND deleted_at IS NULL", teamID, orgID).First(&team).Error; err != nil {
			return err
		}

		var existing models.Membership
		if err := tx.Where("org_id=? AND user_id=? AND team_id=? AND deleted_at IS NULL", orgID, userID, teamID).First(&existing).Error; err == nil {
			return nil
		}

		membership := models.Membership{OrgID: orgID, TeamID: &teamID, OrgUnitID: &teamID, UserID: userID, RoleInTeam: "member"}
		return tx.Create(&membership).Error
	})
}

// ListTeamsForUser returns teams the user belongs to within an org.
func (r *AuthRepo) ListTeamsForUser(ctx context.Context, orgID, userID uuid.UUID) ([]models.Team, error) {
	var teams []models.Team
	err := r.db.WithContext(ctx).
		Table("teams").
		Joins("JOIN memberships m ON m.team_id = teams.id AND m.deleted_at IS NULL").
		Where("teams.org_id = ? AND m.user_id = ? AND teams.deleted_at IS NULL", orgID, userID).
		Find(&teams).Error
	return teams, err
}
