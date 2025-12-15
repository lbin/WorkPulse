package repo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"workpulse/internal/models"
)

type UserAccess struct {
	User        models.User
	OrgUnits    []models.OrgUnit
	Roles       []models.Role
	Permissions []string
}

type RBACRepo struct{ db *gorm.DB }

func NewRBACRepo(db *gorm.DB) *RBACRepo { return &RBACRepo{db: db} }

func (r *RBACRepo) GetUserAccess(ctx context.Context, orgID, userID uuid.UUID) (*UserAccess, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id=? AND org_id=?", userID, orgID).First(&user).Error; err != nil {
		return nil, err
	}

	var memberships []models.Membership
	if err := r.db.WithContext(ctx).Where("org_id=? AND user_id=? AND deleted_at IS NULL", orgID, userID).Find(&memberships).Error; err != nil {
		return nil, err
	}

	roleIDs := make(map[uuid.UUID]struct{})
	unitIDs := make(map[uuid.UUID]struct{})
	for _, m := range memberships {
		if m.RoleID != nil {
			roleIDs[*m.RoleID] = struct{}{}
		}
		if m.OrgUnitID != nil {
			unitIDs[*m.OrgUnitID] = struct{}{}
		} else if m.TeamID != nil {
			unitIDs[*m.TeamID] = struct{}{}
		}
	}

	var roles []models.Role
	if len(roleIDs) > 0 {
		ids := make([]uuid.UUID, 0, len(roleIDs))
		for id := range roleIDs {
			ids = append(ids, id)
		}
		_ = r.db.WithContext(ctx).Where("id IN ?", ids).Find(&roles).Error
	}

	var orgUnits []models.OrgUnit
	if len(unitIDs) > 0 {
		ids := make([]uuid.UUID, 0, len(unitIDs))
		for id := range unitIDs {
			ids = append(ids, id)
		}
		_ = r.db.WithContext(ctx).Where("id IN ?", ids).Find(&orgUnits).Error
	}

	permCodes := r.collectPermissions(ctx, roleIDs)

	return &UserAccess{User: user, OrgUnits: orgUnits, Roles: roles, Permissions: permCodes}, nil
}

func (r *RBACRepo) collectPermissions(ctx context.Context, roleIDs map[uuid.UUID]struct{}) []string {
	if len(roleIDs) == 0 {
		return nil
	}
	ids := make([]uuid.UUID, 0, len(roleIDs))
	for id := range roleIDs {
		ids = append(ids, id)
	}
	var rolePerms []models.RolePermission
	if err := r.db.WithContext(ctx).Where("role_id IN ?", ids).Find(&rolePerms).Error; err != nil {
		return nil
	}
	if len(rolePerms) == 0 {
		return nil
	}
	codes := make([]string, 0, len(rolePerms))
	for _, rp := range rolePerms {
		codes = append(codes, rp.PermissionCode)
	}
	var perms []models.Permission
	if err := r.db.WithContext(ctx).Where("code IN ?", codes).Find(&perms).Error; err != nil {
		return nil
	}
	uniq := make(map[string]struct{})
	for _, p := range perms {
		uniq[p.Code] = struct{}{}
	}
	merged := make([]string, 0, len(uniq))
	for c := range uniq {
		merged = append(merged, c)
	}
	return merged
}
