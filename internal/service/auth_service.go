package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"workpulse/internal/models"
	"workpulse/internal/repo"
	"workpulse/internal/util"
)

// AuthService manages registration and login flows.
type AuthService struct {
	userRepo  *repo.UserRepo
	orgRepo   *repo.OrgRepo
	jwtSecret string
}

// NewAuthService creates a new auth service instance.
func NewAuthService(userRepo *repo.UserRepo, orgRepo *repo.OrgRepo, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, orgRepo: orgRepo, jwtSecret: jwtSecret}
}

// RegisterInput defines the payload for creating a new account.
type RegisterInput struct {
	Email       string
	Password    string
	DisplayName string
	OrgName     string
}

// AuthResult encapsulates information returned after login or signup.
type AuthResult struct {
	Token string
	User  *models.User
}

// Register creates a new organization and primary user.
func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*AuthResult, error) {
	orgName := strings.TrimSpace(in.OrgName)
	if orgName == "" {
		orgName = fmt.Sprintf("%s的团队", strings.Split(in.Email, "@")[0])
	}
	org := &models.Org{Name: orgName, Status: "active"}
	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	hashed, err := util.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		OrgID:        org.ID,
		Email:        strings.ToLower(strings.TrimSpace(in.Email)),
		DisplayName:  strings.TrimSpace(in.DisplayName),
		PasswordHash: hashed,
		Status:       "active",
	}
	if u.DisplayName == "" {
		u.DisplayName = strings.Split(u.Email, "@")[0]
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	token, err := util.NewToken(u.ID.String(), org.ID.String(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, User: u}, nil
}

// Login validates credentials and returns a signed token.
func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	u, err := s.userRepo.FindByEmail(ctx, uuid.Nil, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return nil, err
	}

	if err := util.VerifyPassword(u.PasswordHash, password); err != nil {
		return nil, err
	}

	token, err := util.NewToken(u.ID.String(), u.OrgID.String(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, User: u}, nil
}

// LoginWithOrg restricts login to a specific organization when known.
func (s *AuthService) LoginWithOrg(ctx context.Context, orgID uuid.UUID, email, password string) (*AuthResult, error) {
	u, err := s.userRepo.FindByEmail(ctx, orgID, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return nil, err
	}

	if err := util.VerifyPassword(u.PasswordHash, password); err != nil {
		return nil, err
	}

	token, err := util.NewToken(u.ID.String(), u.OrgID.String(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, User: u}, nil
}

// LoginDefault tries to locate the user without an org hint and falls back to the only record.
func (s *AuthService) LoginDefault(ctx context.Context, email, password string) (*AuthResult, error) {
	return s.Login(ctx, email, password)
}
