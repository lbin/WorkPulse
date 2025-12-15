package service

import (
	"context"

	"github.com/google/uuid"

	"workpulse/internal/models"
	"workpulse/internal/repo"
)

type AuthService struct {
	repo      *repo.RBACRepo
	auditRepo *repo.AuditRepo
}

type UserView struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	AvatarURL   *string   `json:"avatar_url"`
}

type OrgUnitView struct {
	ID           uuid.UUID  `json:"id"`
	OrgID        uuid.UUID  `json:"org_id"`
	ParentUnitID *uuid.UUID `json:"parent_unit_id"`
	Name         string     `json:"name"`
	Path         string     `json:"path"`
	UnitType     string     `json:"unit_type"`
}

type RoleView struct {
	ID        uuid.UUID  `json:"id"`
	OrgID     uuid.UUID  `json:"org_id"`
	OrgUnitID *uuid.UUID `json:"org_unit_id"`
	Name      string     `json:"name"`
	Desc      *string    `json:"description"`
}

type AuthContext struct {
	User          UserView      `json:"user"`
	OrgUnits      []OrgUnitView `json:"org_units"`
	Roles         []RoleView    `json:"roles"`
	Permissions   []string      `json:"permissions"`
	ActiveOrgUnit *uuid.UUID    `json:"active_org_unit"`
}

func NewAuthService(repo *repo.RBACRepo, auditRepo *repo.AuditRepo) *AuthService {
	return &AuthService{repo: repo, auditRepo: auditRepo}
}

func (s *AuthService) GetContext(ctx context.Context, orgID, userID uuid.UUID, orgUnitID *uuid.UUID) (*AuthContext, error) {
	access, err := s.repo.GetUserAccess(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}
	perms := mergePermissions(access.Permissions, defaultPermissions())
	active := resolveOrgUnit(orgUnitID, access.OrgUnits)
	return &AuthContext{
		User:          toUserView(access.User),
		OrgUnits:      toOrgUnitViews(access.OrgUnits),
		Roles:         toRoleViews(access.Roles),
		Permissions:   perms,
		ActiveOrgUnit: active,
	}, nil
}

func mergePermissions(src []string, defaults []string) []string {
	set := make(map[string]struct{})
	for _, p := range defaults {
		set[p] = struct{}{}
	}
	for _, p := range src {
		set[p] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	return out
}

func resolveOrgUnit(requested *uuid.UUID, units []models.OrgUnit) *uuid.UUID {
	if requested != nil {
		for _, u := range units {
			if u.ID == *requested {
				return requested
			}
		}
	}
	if len(units) == 0 {
		return nil
	}
	return &units[0].ID
}

func defaultPermissions() []string {
	return []string{
		"dashboard.view",
		"projects.view",
		"projects.manage",
		"okr.view",
		"okr.manage",
		"meetings.view",
		"meetings.manage",
		"reports.view",
		"reports.manage",
		"reports.review",
		"workitems.manage",
		"analytics.view",
	}
}

func toUserView(u models.User) UserView {
	return UserView{ID: u.ID, OrgID: u.OrgID, Email: u.Email, DisplayName: u.DisplayName, AvatarURL: u.AvatarURL}
}

func toOrgUnitViews(items []models.OrgUnit) []OrgUnitView {
	res := make([]OrgUnitView, 0, len(items))
	for _, it := range items {
		res = append(res, OrgUnitView{ID: it.ID, OrgID: it.OrgID, ParentUnitID: it.ParentUnitID, Name: it.Name, Path: it.Path, UnitType: it.UnitType})
	}
	return res
}

func toRoleViews(items []models.Role) []RoleView {
	res := make([]RoleView, 0, len(items))
	for _, it := range items {
		res = append(res, RoleView{ID: it.ID, OrgID: it.OrgID, OrgUnitID: it.OrgUnitID, Name: it.Name, Desc: it.Description})
	}
	return res
}
