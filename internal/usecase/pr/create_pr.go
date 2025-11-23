package pr

import (
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/interfaces"
)

type CreatePRUseCase struct {
	prRepo   interfaces.PRRepository
	userRepo interfaces.UserRepository
	teamRepo interfaces.TeamRepository
	reviewer *domain.ReviewerAssigner
}

func NewCreatePRUseCase(
	prRepo interfaces.PRRepository,
	userRepo interfaces.UserRepository,
	teamRepo interfaces.TeamRepository,
	reviewer *domain.ReviewerAssigner,
) *CreatePRUseCase {
	return &CreatePRUseCase{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
		reviewer: reviewer,
	}
}

type CreatePRRequest struct {
	PRID     string
	PRName   string
	AuthorID string
}

func (uc *CreatePRUseCase) Execute(req CreatePRRequest) (*domain.PullRequest, error) {
	exists, err := uc.prRepo.Exists(req.PRID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrPRExists
	}

	author, err := uc.userRepo.GetByID(req.AuthorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, domain.ErrNotFound
	}

	team, err := uc.teamRepo.GetByName(author.TeamName)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, domain.ErrNotFound
	}

	reviewers := uc.reviewer.AssignReviewers(team, author, 2)

	pr := domain.NewPullRequest(req.PRID, req.PRName, req.AuthorID, reviewers)

	if err := uc.prRepo.Create(pr); err != nil {
		return nil, err
	}

	return pr, nil
}
