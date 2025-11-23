package domain

import "errors"

var (
	ErrTeamExists      = errors.New("TEAM_EXISTS")
	ErrPRExists        = errors.New("PR_EXISTS")
	ErrPRMerged        = errors.New("PR_MERGED")
	ErrNotAssigned     = errors.New("NOT_ASSIGNED")
	ErrNoCandidate     = errors.New("NO_CANDIDATE")
	ErrNotFound        = errors.New("NOT_FOUND")
	ErrInvalidStatus   = errors.New("INVALID_STATUS")
	ErrInvalidArgument = errors.New("INVALID_ARGUMENT")
)

type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     errors.New(code),
	}
}
