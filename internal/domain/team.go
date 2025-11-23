package domain

type Team struct {
	TeamName string
	Members  []*User
}

func NewTeam(teamName string, members []*User) *Team {
	return &Team{
		TeamName: teamName,
		Members:  members,
	}
}

func (t *Team) GetActiveMembers() []*User {
	var active []*User
	for _, member := range t.Members {
		if member.IsActive {
			active = append(active, member)
		}
	}

	return active
}

func (t *Team) GetActiveMembersExcluding(excludeUserID string) []*User {
	var active []*User
	for _, member := range t.Members {
		if member.IsActive && member.UserID != excludeUserID {
			active = append(active, member)
		}
	}

	return active
}
