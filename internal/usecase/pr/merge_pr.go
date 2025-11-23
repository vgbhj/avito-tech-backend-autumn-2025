package pr

import (
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/interfaces"
)

type MergePRUseCase struct {
	prRepo interfaces.PRRepository
}

func NewMergePRUseCase(prRepo interfaces.PRRepository) *MergePRUseCase {
	return &MergePRUseCase{
		prRepo: prRepo,
	}
}

type MergePRRequest struct {
	PRID string
}

func (uc *MergePRUseCase) Execute(req MergePRRequest) (*domain.PullRequest, error) {
	pr, err := uc.prRepo.GetByID(req.PRID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, domain.ErrNotFound
	}

	if err := pr.Merge(); err != nil {
		return nil, err
	}

	if err := uc.prRepo.Update(pr); err != nil {
		return nil, err
	}

	return pr, nil
}
