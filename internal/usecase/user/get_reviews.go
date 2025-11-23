package user

import (
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/interfaces"
)

type GetReviewsUseCase struct {
	prRepo   interfaces.PRRepository
	userRepo interfaces.UserRepository
}

func NewGetReviewsUseCase(prRepo interfaces.PRRepository, userRepo interfaces.UserRepository) *GetReviewsUseCase {
	return &GetReviewsUseCase{
		prRepo:   prRepo,
		userRepo: userRepo,
	}
}

func (uc *GetReviewsUseCase) Execute(userID string) ([]*domain.PullRequest, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domain.ErrNotFound
	}

	prs, err := uc.prRepo.GetByReviewerID(userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
