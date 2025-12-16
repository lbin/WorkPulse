package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"workpulse/internal/models"
	"workpulse/internal/repo"
)

type AuthService struct {
	repo      *repo.RBACRepo
	auditRepo *repo.AuditRepo
	authRepo  *repo.AuthRepo
	jwtSecret string
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

// NewAuthService wires up the auth service with its dependencies.
func NewAuthService(repo *repo.RBACRepo, auditRepo *repo.AuditRepo, authRepo *repo.AuthRepo, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, auditRepo: auditRepo, authRepo: authRepo, jwtSecret: jwtSecret}
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

// Login registers a session token for an existing user based on credentials.
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := s.authRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.issueToken(user.ID, user.OrgID)
	if err != nil {
		return nil, err
	}

	ctxView, err := s.GetContext(ctx, user.OrgID, user.ID, nil)
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: token, Context: ctxView}, nil
}

// Register creates a brand-new org and user using email/password credentials.
func (s *AuthService) Register(ctx context.Context, email, password, displayName string) (*LoginResult, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	if displayName == "" {
		parts := strings.Split(email, "@")
		if len(parts) > 0 {
			displayName = parts[0]
		}
	}

	if _, err := s.authRepo.FindUserByEmail(ctx, email); err == nil {
		return nil, errors.New("account already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, _, err := s.authRepo.CreateUserWithOrg(ctx, email, displayName, string(hash))
	if err != nil {
		return nil, err
	}

	if _, err := s.authRepo.CreateTeamWithMembership(ctx, user.OrgID, user.ID, "Default Team", nil); err != nil {
		return nil, err
	}

	token, err := s.issueToken(user.ID, user.OrgID)
	if err != nil {
		return nil, err
	}

	ctxView, err := s.GetContext(ctx, user.OrgID, user.ID, nil)
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: token, Context: ctxView}, nil
}

// CreateTeam provisions a new team/org unit and enrolls the requesting user.
func (s *AuthService) CreateTeam(ctx context.Context, orgID, userID uuid.UUID, name string, parentTeamID *uuid.UUID) (*models.Team, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("team name is required")
	}
	return s.authRepo.CreateTeamWithMembership(ctx, orgID, userID, name, parentTeamID)
}

// JoinTeam subscribes a user to an existing team in their org.
func (s *AuthService) JoinTeam(ctx context.Context, orgID, userID, teamID uuid.UUID) error {
	return s.authRepo.AddUserToTeam(ctx, orgID, userID, teamID)
}

// ListTeams shows a user's memberships inside an org.
func (s *AuthService) ListTeams(ctx context.Context, orgID, userID uuid.UUID) ([]models.Team, error) {
	return s.authRepo.ListTeamsForUser(ctx, orgID, userID)
}

// issueToken builds a signed JWT with org and user claims.
func (s *AuthService) issueToken(userID, orgID uuid.UUID) (string, error) {
	type authClaims struct {
		UserID string `json:"user_id"`
		OrgID  string `json:"org_id"`
		jwt.RegisteredClaims
	}
	claims := authClaims{
		UserID:           userID.String(),
		OrgID:            orgID.String(),
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// LoginResult returns the JWT token and hydrated context.
type LoginResult struct {
	Token   string       `json:"token"`
	Context *AuthContext `json:"context"`
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
		"teams.manage",
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
