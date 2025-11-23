package interfaces

import "github.com/avito-tech-backend-autumn-2025/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error

	Update(user *domain.User) error

	GetByID(userID string) (*domain.User, error)

	GetByTeamName(teamName string) ([]*domain.User, error)

	Exists(userID string) (bool, error)
}
