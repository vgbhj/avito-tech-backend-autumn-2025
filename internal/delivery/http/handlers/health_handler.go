package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
// @Summary      Health check
// @Description  Проверка работоспособности сервиса
// @Tags         Health
// @Accept       json
// @Produce      text/plain
// @Success      200  {string}  string  "OK"
// @Router       /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func (h *HealthHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)
}
