package dto

import (
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/team"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/user"
)

func ToTeamDTO(team *domain.Team) TeamDTO {
	members := make([]TeamMemberDTO, 0, len(team.Members))
	for _, member := range team.Members {
		members = append(members, TeamMemberDTO{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}
	return TeamDTO{
		TeamName: team.TeamName,
		Members:  members,
	}
}

func ToUserDTO(user *domain.User) UserDTO {
	return UserDTO{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func ToPullRequestDTO(pr *domain.PullRequest) PullRequestDTO {
	return PullRequestDTO{
		PRID:              pr.ID,
		PRName:            pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func ToPullRequestShortDTO(pr *domain.PullRequest) PullRequestShortDTO {
	return PullRequestShortDTO{
		PRID:     pr.ID,
		PRName:   pr.Name,
		AuthorID: pr.AuthorID,
		Status:   string(pr.Status),
	}
}

func ToCreateTeamRequest(req CreateTeamRequest) team.CreateTeamRequest {
	members := make([]team.TeamMemberRequest, 0, len(req.Members))
	for _, member := range req.Members {
		members = append(members, team.TeamMemberRequest{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}
	return team.CreateTeamRequest{
		TeamName: req.TeamName,
		Members:  members,
	}
}

func ToSetActiveRequest(req SetActiveRequest) user.SetActiveRequest {
	return user.SetActiveRequest{
		UserID:   req.UserID,
		IsActive: req.IsActive,
	}
}
