package interfaces

import "github.com/avito-tech-backend-autumn-2025/internal/domain"

type TeamRepository interface {
	Create(team *domain.Team) error
	GetByName(teamName string) (*domain.Team, error)
	Exists(teamNmae string) (bool, error)
}
