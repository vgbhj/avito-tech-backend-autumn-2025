package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/avito-tech-backend-autumn-2025/internal/delivery/http/dto"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/pr"
)

type PRHandler struct {
	createPRUseCase         *pr.CreatePRUseCase
	mergePRUseCase          *pr.MergePRUseCase
	reassignReviewerUseCase *pr.ReassignReviewerUseCase
}

func NewPRHandler(
	createPRUseCase *pr.CreatePRUseCase,
	mergePRUseCase *pr.MergePRUseCase,
	reassignReviewerUseCase *pr.ReassignReviewerUseCase,
) *PRHandler {
	return &PRHandler{
		createPRUseCase:         createPRUseCase,
		mergePRUseCase:          mergePRUseCase,
		reassignReviewerUseCase: reassignReviewerUseCase,
	}
}

// CreatePR godoc
// @Summary      Создать PR и назначить ревьюеров
// @Description  Создаёт PR и автоматически назначает до 2 ревьюеров из команды автора
// @Tags         PullRequests
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreatePRRequest  true  "Данные PR"
// @Success      201      {object}  dto.PRResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      409      {object}  dto.ErrorResponse
// @Router       /pullRequest/create [post]
func (h *PRHandler) CreatePR(c *gin.Context) {
	var req dto.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	useCaseReq := pr.CreatePRRequest{
		PRID:     req.PRID,
		PRName:   req.PRName,
		AuthorID: req.AuthorID,
	}

	pr, err := h.createPRUseCase.Execute(useCaseReq)
	if err != nil {
		handleDomainError(c, err)
		return
	}

	response := dto.PRResponse{
		PR: dto.ToPullRequestDTO(pr),
	}

	respondJSON(c, http.StatusCreated, response)
}

// MergePR godoc
// @Summary      Пометить PR как MERGED
// @Description  Помечает PR как MERGED (идемпотентная операция)
// @Tags         PullRequests
// @Accept       json
// @Produce      json
// @Param        request  body      dto.MergePRRequest  true  "ID PR"
// @Success      200      {object}  dto.PRResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Router       /pullRequest/merge [post]
func (h *PRHandler) MergePR(c *gin.Context) {
	var req dto.MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	useCaseReq := pr.MergePRRequest{
		PRID: req.PRID,
	}

	pr, err := h.mergePRUseCase.Execute(useCaseReq)
	if err != nil {
		handleDomainError(c, err)
		return
	}

	response := dto.PRResponse{
		PR: dto.ToPullRequestDTO(pr),
	}

	respondJSON(c, http.StatusOK, response)
}

// ReassignReviewer godoc
// @Summary      Переназначить ревьюера
// @Description  Переназначает конкретного ревьювера на другого из его команды
// @Tags         PullRequests
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ReassignReviewerRequest  true  "Данные переназначения"
// @Success      200      {object}  dto.ReassignReviewerResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      409      {object}  dto.ErrorResponse
// @Router       /pullRequest/reassign [post]
func (h *PRHandler) ReassignReviewer(c *gin.Context) {
	var req dto.ReassignReviewerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	useCaseReq := pr.ReassignReviewerRequest{
		PRID:      req.PRID,
		OldUserID: req.OldUserID,
	}

	result, err := h.reassignReviewerUseCase.Execute(useCaseReq)
	if err != nil {
		handleDomainError(c, err)
		return
	}

	response := dto.ReassignReviewerResponse{
		PR:         dto.ToPullRequestDTO(result.PR),
		ReplacedBy: result.ReplacedBy,
	}

	respondJSON(c, http.StatusOK, response)
}

func (h *PRHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/pullRequest/create", h.CreatePR)
	r.POST("/pullRequest/merge", h.MergePR)
	r.POST("/pullRequest/reassign", h.ReassignReviewer)
}
