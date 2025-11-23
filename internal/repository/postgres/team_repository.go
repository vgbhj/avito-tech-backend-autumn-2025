package postgres

import (
	"database/sql"

	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/interfaces"
)

// teamRepository реализует TeamRepository для PostgreSQL
type teamRepository struct {
	db *sql.DB
}

// NewTeamRepository создает новый teamRepository
func NewTeamRepository(db *sql.DB) interfaces.TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(team *domain.Team) error {
	query := `INSERT INTO teams (team_name, created_at, updated_at) 
	          VALUES ($1, NOW(), NOW()) 
	          ON CONFLICT (team_name) DO NOTHING`

	_, err := r.db.Exec(query, team.TeamName)
	if err != nil {
		return err
	}

	return nil
}

func (r *teamRepository) GetByName(teamName string) (*domain.Team, error) {
	var team domain.Team
	query := `SELECT team_name FROM teams WHERE team_name = $1`
	err := r.db.QueryRow(query, teamName).Scan(&team.TeamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	members, err := r.getTeamMembers(teamName)
	if err != nil {
		return nil, err
	}
	team.Members = members
	return &team, nil
}

func (r *teamRepository) getTeamMembers(teamName string) ([]*domain.User, error) {
	query := `SELECT user_id, username, team_name, is_active 
	          FROM users 
	          WHERE team_name = $1`

	rows, err := r.db.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, err
		}
		members = append(members, &user)
	}

	return members, rows.Err()
}

func (r *teamRepository) Exists(teamName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`
	err := r.db.QueryRow(query, teamName).Scan(&exists)
	return exists, err
}
