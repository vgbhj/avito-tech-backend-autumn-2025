package domain

import (
	"math/rand"
	"time"
)

type ReviewerAssigner struct{}

func NewReviewerAssigner() *ReviewerAssigner {
	return &ReviewerAssigner{}
}

func (ra *ReviewerAssigner) AssignReviewers(team *Team, author *User, maxReviewers int) []string {
	if maxReviewers <= 0 {
		return []string{}
	}

	candidates := team.GetActiveMembersExcluding(author.UserID)

	if len(candidates) == 0 {
		return []string{}
	}

	count := len(candidates)
	if count > maxReviewers {
		count = maxReviewers
	}

	shuffled := ra.shuffle(candidates)

	reviewers := make([]string, 0, count)
	for i := 0; i < count; i++ {
		reviewers = append(reviewers, shuffled[i].UserID)
	}

	return reviewers
}

func (ra *ReviewerAssigner) FindReplacementCandidate(team *Team, excludeUserIDs []string) (*User, error) {
	excludeMap := make(map[string]bool)
	for _, id := range excludeUserIDs {
		excludeMap[id] = true
	}

	var candidates []*User
	for _, member := range team.GetActiveMembers() {
		if !excludeMap[member.UserID] {
			candidates = append(candidates, member)
		}
	}

	if len(candidates) == 0 {
		return nil, ErrNoCandidate
	}

	shuffled := ra.shuffle(candidates)
	return shuffled[0], nil
}

func (ra *ReviewerAssigner) shuffle(users []*User) []*User {
	shuffled := make([]*User, len(users))
	copy(shuffled, users)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}
