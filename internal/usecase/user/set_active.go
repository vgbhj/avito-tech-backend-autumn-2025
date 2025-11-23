package user

import (
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/interfaces"
)

type SetActiveUseCase struct {
	userRepo interfaces.UserRepository
}

func NewSetActiveUseCase(userRepo interfaces.UserRepository) *SetActiveUseCase {
	return &SetActiveUseCase{
		userRepo: userRepo,
	}
}

type SetActiveRequest struct {
	UserID   string
	IsActive bool
}

func (uc *SetActiveUseCase) Execute(req SetActiveRequest) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domain.ErrNotFound
	}

	user.SetActive(req.IsActive)

	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
