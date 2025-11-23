package interfaces

import "github.com/avito-tech-backend-autumn-2025/internal/domain"

type PRRepository interface {
	Create(pr *domain.PullRequest) error

	Update(pr *domain.PullRequest) error

	GetById(prID string) (*domain.PullRequest, error)

	GetByReviewerID(reviewerID string) ([]*domain.PullRequest, error)

	Exists(prID string) (bool, error)
}
