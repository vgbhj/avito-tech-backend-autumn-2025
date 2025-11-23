package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/avito-tech-backend-autumn-2025/internal/delivery/http/dto"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/team"
)

type TeamHandler struct {
	createTeamUseCase *team.CreateTeamUseCase
	getTeamUseCase    *team.GetTeamUseCase
}

func NewTeamHandler(
	createTeamUseCase *team.CreateTeamUseCase,
	getTeamUseCase *team.GetTeamUseCase,
) *TeamHandler {
	return &TeamHandler{
		createTeamUseCase: createTeamUseCase,
		getTeamUseCase:    getTeamUseCase,
	}
}

// CreateTeam godoc
// @Summary      Создать команду с участниками
// @Description  Создаёт команду с участниками (создаёт/обновляет пользователей)
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Param        team  body      dto.CreateTeamRequest  true  "Данные команды"
// @Success      201   {object}  dto.TeamResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Router       /team/add [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	useCaseReq := dto.ToCreateTeamRequest(req)
	team, err := h.createTeamUseCase.Execute(useCaseReq)
	if err != nil {
		handleDomainError(c, err)
		return
	}

	response := dto.TeamResponse{
		Team: dto.ToTeamDTO(team),
	}

	respondJSON(c, http.StatusCreated, response)
}

// GetTeam godoc
// @Summary      Получить команду с участниками
// @Description  Получает команду с участниками по имени
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Param        team_name  query     string  true  "Уникальное имя команды"
// @Success      200        {object}  dto.TeamDTO
// @Failure      404        {object}  dto.ErrorResponse
// @Router       /team/get [get]
func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "team_name is required")
		return
	}

	team, err := h.getTeamUseCase.Execute(teamName)
	if err != nil {
		handleDomainError(c, err)
		return
	}

	respondJSON(c, http.StatusOK, dto.ToTeamDTO(team))
}

func (h *TeamHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/team/add", h.CreateTeam)
	r.GET("/team/get", h.GetTeam)
}
