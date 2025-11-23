package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/avito-tech-backend-autumn-2025/internal/delivery/http/dto"
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
)

func respondJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

func respondError(c *gin.Context, statusCode int, code, message string) {
	response := dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	c.JSON(statusCode, response)
}

func handleDomainError(c *gin.Context, err error) {
	switch err {
	case domain.ErrTeamExists:
		respondError(c, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
	case domain.ErrPRExists:
		respondError(c, http.StatusConflict, "PR_EXISTS", "PR id already exists")
	case domain.ErrPRMerged:
		respondError(c, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
	case domain.ErrNotAssigned:
		respondError(c, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
	case domain.ErrNoCandidate:
		respondError(c, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
	case domain.ErrNotFound:
		respondError(c, http.StatusNotFound, "NOT_FOUND", "resource not found")
	default:
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
