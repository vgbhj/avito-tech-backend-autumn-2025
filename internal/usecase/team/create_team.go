package team

import (
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/interfaces"
)

type CreateTeamUseCase struct {
	teamRepo interfaces.TeamRepository
	userRepo interfaces.UserRepository
}

func NewCreateTeamUseCase(teamRepo interfaces.TeamRepository, userRepo interfaces.UserRepository) *CreateTeamUseCase {
	return &CreateTeamUseCase{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

type CreateTeamRequest struct {
	TeamName string
	Members  []TeamMemberRequest
}

type TeamMemberRequest struct {
	UserID   string
	Username string
	IsActive bool
}

func (uc *CreateTeamUseCase) Execute(req CreateTeamRequest) (*domain.Team, error) {
	exists, err := uc.teamRepo.Exists(req.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrTeamExists
	}

	team := domain.NewTeam(req.TeamName, nil)
	if err := uc.teamRepo.Create(team); err != nil {
		return nil, err
	}

	users := make([]*domain.User, 0, len(req.Members))
	for _, memberReq := range req.Members {
		userExists, err := uc.userRepo.Exists(memberReq.UserID)
		if err != nil {
			return nil, err
		}

		user := domain.NewUser(memberReq.UserID, memberReq.Username, req.TeamName, memberReq.IsActive)

		if userExists {
			user, err = uc.userRepo.GetByID(memberReq.UserID)
			if err != nil {
				return nil, err
			}
			user.TeamName = req.TeamName
			user.IsActive = memberReq.IsActive
			if err := uc.userRepo.Update(user); err != nil {
				return nil, err
			}
		} else {
			if err := uc.userRepo.Create(user); err != nil {
				return nil, err
			}
		}

		users = append(users, user)
	}

	team.Members = users
	return team, nil
}
