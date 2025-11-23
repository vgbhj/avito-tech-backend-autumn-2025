package postgres

import (
	"database/sql"
	"time"

	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/interfaces"
)

type prRepository struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) interfaces.PRRepository {
	return &prRepository{db: db}
}

func (r *prRepository) Create(pr *domain.PullRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	var mergedAt *time.Time
	if pr.MergedAt != nil {
		mergedAt = pr.MergedAt
	}

	_, err = tx.Exec(query, pr.ID, pr.Name, pr.AuthorID, string(pr.Status), pr.CreatedAt, mergedAt)
	if err != nil {
		return err
	}

	for _, reviewerID := range pr.AssignedReviewers {
		reviewerQuery := `INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) 
		                  VALUES ($1, $2, NOW())`
		if _, err := tx.Exec(reviewerQuery, pr.ID, reviewerID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *prRepository) Update(pr *domain.PullRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE pull_requests 
	          SET pull_request_name = $2, status = $3, merged_at = $4 
	          WHERE pull_request_id = $1`

	var mergedAt *time.Time
	if pr.MergedAt != nil {
		mergedAt = pr.MergedAt
	}

	_, err = tx.Exec(query, pr.ID, pr.Name, string(pr.Status), mergedAt)
	if err != nil {
		return err
	}

	deleteQuery := `DELETE FROM pr_reviewers WHERE pull_request_id = $1`
	if _, err := tx.Exec(deleteQuery, pr.ID); err != nil {
		return err
	}

	for _, reviewerID := range pr.AssignedReviewers {
		reviewerQuery := `INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) 
		                  VALUES ($1, $2, NOW())`
		if _, err := tx.Exec(reviewerQuery, pr.ID, reviewerID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *prRepository) GetByID(prID string) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	var statusStr string
	var mergedAt sql.NullTime

	query := `SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at 
	          FROM pull_requests 
	          WHERE pull_request_id = $1`

	err := r.db.QueryRow(query, prID).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &statusStr, &pr.CreatedAt, &mergedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	pr.Status = domain.PRStatus(statusStr)
	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	reviewers, err := r.getReviewers(prID)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *prRepository) getReviewers(prID string) ([]string, error) {
	query := `SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $1 ORDER BY assigned_at`

	rows, err := r.db.Query(query, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, reviewerID)
	}

	return reviewers, rows.Err()
}

func (r *prRepository) GetByReviewerID(reviewerID string) ([]*domain.PullRequest, error) {
	query := `SELECT DISTINCT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at
	          FROM pull_requests pr
	          INNER JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
	          WHERE prr.reviewer_id = $1
	          ORDER BY pr.created_at DESC`

	rows, err := r.db.Query(query, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		var statusStr string
		var mergedAt sql.NullTime

		if err := rows.Scan(
			&pr.ID, &pr.Name, &pr.AuthorID, &statusStr, &pr.CreatedAt, &mergedAt,
		); err != nil {
			return nil, err
		}

		pr.Status = domain.PRStatus(statusStr)
		if mergedAt.Valid {
			pr.MergedAt = &mergedAt.Time
		}

		reviewers, err := r.getReviewers(pr.ID)
		if err != nil {
			return nil, err
		}
		pr.AssignedReviewers = reviewers

		prs = append(prs, &pr)
	}

	return prs, rows.Err()
}

func (r *prRepository) Exists(prID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`
	err := r.db.QueryRow(query, prID).Scan(&exists)
	return exists, err
}
