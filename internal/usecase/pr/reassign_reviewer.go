package pr

import (
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/interfaces"
)

type ReassignReviewerUseCase struct {
	prRepo   interfaces.PRRepository
	userRepo interfaces.UserRepository
	teamRepo interfaces.TeamRepository
	reviewer *domain.ReviewerAssigner
}

func NewReassignReviewerUseCase(
	prRepo interfaces.PRRepository,
	userRepo interfaces.UserRepository,
	teamRepo interfaces.TeamRepository,
	reviewer *domain.ReviewerAssigner,
) *ReassignReviewerUseCase {
	return &ReassignReviewerUseCase{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
		reviewer: reviewer,
	}
}

type ReassignReviewerRequest struct {
	PRID      string
	OldUserID string
}

type ReassignReviewerResponse struct {
	PR         *domain.PullRequest
	ReplacedBy string
}

func (uc *ReassignReviewerUseCase) Execute(req ReassignReviewerRequest) (*ReassignReviewerResponse, error) {
	pr, err := uc.prRepo.GetByID(req.PRID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, domain.ErrNotFound
	}

	if !pr.CanReassign() {
		return nil, domain.ErrPRMerged
	}

	if !pr.HasReviewer(req.OldUserID) {
		return nil, domain.ErrNotAssigned
	}

	oldReviewer, err := uc.userRepo.GetByID(req.OldUserID)
	if err != nil {
		return nil, err
	}
	if oldReviewer == nil {
		return nil, domain.ErrNotFound
	}

	team, err := uc.teamRepo.GetByName(oldReviewer.TeamName)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, domain.ErrNotFound
	}

	excludeIDs := []string{pr.AuthorID}
	excludeIDs = append(excludeIDs, pr.AssignedReviewers...)

	newReviewer, err := uc.reviewer.FindReplacementCandidate(team, excludeIDs)
	if err != nil {
		return nil, err
	}

	if err := pr.ReplaceReviewer(req.OldUserID, newReviewer.UserID); err != nil {
		return nil, err
	}

	if err := uc.prRepo.Update(pr); err != nil {
		return nil, err
	}

	return &ReassignReviewerResponse{
		PR:         pr,
		ReplacedBy: newReviewer.UserID,
	}, nil
}
