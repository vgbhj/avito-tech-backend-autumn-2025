package domain

import "time"

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            PRStatus
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          *time.Time
}

func NewPullRequest(id, name, authorID string, reviewers []string) *PullRequest {
	return &PullRequest{
		ID:                id,
		Name:              name,
		AuthorID:          authorID,
		Status:            StatusOpen,
		AssignedReviewers: reviewers,
		CreatedAt:         time.Now(),
		MergedAt:          nil,
	}
}

func (pr *PullRequest) Merge() error {
	if pr.Status == StatusMerged {
		return nil
	}

	pr.Status = StatusMerged
	now := time.Now()
	pr.MergedAt = &now
	return nil
}

func (pr *PullRequest) CanReassign() bool {
	return pr.Status == StatusOpen
}

func (pr *PullRequest) IsMerged() bool {
	return pr.Status == StatusMerged
}

func (pr *PullRequest) HasReviewer(userID string) bool {
	for _, reviewerID := range pr.AssignedReviewers {
		if reviewerID == userID {
			return true
		}
	}
	return false
}

func (pr *PullRequest) ReplaceReviewer(oldUserID, newUserID string) error {
	if !pr.CanReassign() {
		return ErrPRMerged
	}

	if !pr.HasReviewer(oldUserID) {
		return ErrNotAssigned
	}

	for i, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldUserID {
			pr.AssignedReviewers[i] = newUserID
			return nil
		}
	}

	return ErrNotAssigned
}
