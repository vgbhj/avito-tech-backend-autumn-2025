package domain

type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}

func NewUser(userID, username, teamName string, isActive bool) *User {
	return &User{
		UserID:   userID,
		Username: username,
		TeamName: teamName,
		IsActive: isActive,
	}
}

func (u *User) SetActive(isActive bool) {
	u.IsActive = isActive
}
