package team

import (
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/interfaces"
)

type GetTeamUseCase struct {
	teamRepo interfaces.TeamRepository
}

func NewGetTeamUseCase(teamRepo interfaces.TeamRepository) *GetTeamUseCase {
	return &GetTeamUseCase{
		teamRepo: teamRepo,
	}
}

func (uc *GetTeamUseCase) Execute(teamName string) (*domain.Team, error) {
	team, err := uc.teamRepo.GetByName(teamName)
	if err != nil {
		return nil, err
	}

	if team == nil {
		return nil, domain.ErrNotFound
	}

	return team, nil
}
