package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/avito-tech-backend-autumn-2025/internal/delivery/http/dto"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/user"
)

type UserHandler struct {
	setActiveUseCase  *user.SetActiveUseCase
	getReviewsUseCase *user.GetReviewsUseCase
}

func NewUserHandler(
	setActiveUseCase *user.SetActiveUseCase,
	getReviewsUseCase *user.GetReviewsUseCase,
) *UserHandler {
	return &UserHandler{
		setActiveUseCase:  setActiveUseCase,
		getReviewsUseCase: getReviewsUseCase,
	}
}

// SetActive godoc
// @Summary      Установить флаг активности пользователя
// @Description  Устанавливает флаг активности пользователя
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SetActiveRequest  true  "Данные пользователя"
// @Success      200      {object}  dto.UserResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Router       /users/setIsActive [post]
func (h *UserHandler) SetActive(c *gin.Context) {
	var req dto.SetActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	useCaseReq := dto.ToSetActiveRequest(req)
	user, err := h.setActiveUseCase.Execute(useCaseReq)
	if err != nil {
		handleDomainError(c, err)
		return
	}

	response := dto.UserResponse{
		User: dto.ToUserDTO(user),
	}

	respondJSON(c, http.StatusOK, response)
}

// GetReviews godoc
// @Summary      Получить PR'ы пользователя
// @Description  Получает PR'ы, где пользователь назначен ревьюером
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id  query     string  true  "Идентификатор пользователя"
// @Success      200      {object}  dto.GetReviewsResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Router       /users/getReview [get]
func (h *UserHandler) GetReviews(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "user_id is required")
		return
	}

	prs, err := h.getReviewsUseCase.Execute(userID)
	if err != nil {
		handleDomainError(c, err)
		return
	}

	prDTOs := make([]dto.PullRequestShortDTO, 0, len(prs))
	for _, pr := range prs {
		prDTOs = append(prDTOs, dto.ToPullRequestShortDTO(pr))
	}

	response := dto.GetReviewsResponse{
		UserID:       userID,
		PullRequests: prDTOs,
	}

	respondJSON(c, http.StatusOK, response)
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/users/setIsActive", h.SetActive)
	r.GET("/users/getReview", h.GetReviews)
}
