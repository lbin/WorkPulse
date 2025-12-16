package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"workpulse/internal/models"
	"workpulse/internal/repo"
)

// TeamService wraps team and membership operations.
type TeamService struct {
	teamRepo   *repo.TeamRepo
	memberRepo *repo.MembershipRepo
}

// NewTeamService constructs a new TeamService instance.
func NewTeamService(teamRepo *repo.TeamRepo, memberRepo *repo.MembershipRepo) *TeamService {
	return &TeamService{teamRepo: teamRepo, memberRepo: memberRepo}
}

// CreateTeam builds a new team for the org and attaches the creator as owner.
func (s *TeamService) CreateTeam(ctx context.Context, orgID, creatorID uuid.UUID, name string) (*models.Team, error) {
	cleaned := strings.TrimSpace(name)
	if cleaned == "" {
		return nil, fmt.Errorf("team name required")
	}
	team := &models.Team{OrgID: orgID, Name: cleaned, Path: cleaned}
	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}

	m := &models.Membership{OrgID: orgID, TeamID: team.ID, UserID: creatorID, RoleInTeam: "owner"}
	if err := s.memberRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return team, nil
}

// ListTeams enumerates all teams inside the organization.
func (s *TeamService) ListTeams(ctx context.Context, orgID uuid.UUID) ([]models.Team, error) {
	return s.teamRepo.List(ctx, orgID)
}

// ListUserTeams returns all memberships for the user.
func (s *TeamService) ListUserTeams(ctx context.Context, userID uuid.UUID) ([]models.Membership, error) {
	return s.memberRepo.ListByUser(ctx, userID)
}

// JoinTeam creates a membership if it does not exist.
func (s *TeamService) JoinTeam(ctx context.Context, orgID, teamID, userID uuid.UUID) error {
	team, err := s.teamRepo.Get(ctx, orgID, teamID)
	if err != nil {
		return err
	}
	exists, err := s.memberRepo.Exists(ctx, team.ID, userID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	m := &models.Membership{OrgID: orgID, TeamID: team.ID, UserID: userID, RoleInTeam: "member"}
	return s.memberRepo.Create(ctx, m)
}
