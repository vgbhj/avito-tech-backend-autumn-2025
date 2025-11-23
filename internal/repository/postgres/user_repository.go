package postgres

import (
	"database/sql"

	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/interfaces"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (user_id, username, team_name, is_active, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, NOW(), NOW())`

	_, err := r.db.Exec(query, user.UserID, user.Username, user.TeamName, user.IsActive)
	return err
}

func (r *userRepository) Update(user *domain.User) error {
	query := `UPDATE users 
	          SET username = $2, team_name = $3, is_active = $4, updated_at = NOW() 
	          WHERE user_id = $1`

	_, err := r.db.Exec(query, user.UserID, user.Username, user.TeamName, user.IsActive)
	return err
}

func (r *userRepository) GetByID(userID string) (*domain.User, error) {
	var user domain.User
	query := `SELECT user_id, username, team_name, is_active 
	          FROM users 
	          WHERE user_id = $1`

	err := r.db.QueryRow(query, userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByTeamName(teamName string) ([]*domain.User, error) {
	query := `SELECT user_id, username, team_name, is_active 
	          FROM users 
	          WHERE team_name = $1`

	rows, err := r.db.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, rows.Err()
}

func (r *userRepository) Exists(userID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`
	err := r.db.QueryRow(query, userID).Scan(&exists)
	return exists, err
}
