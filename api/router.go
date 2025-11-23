package api

import (
	"github.com/gin-gonic/gin"

	"github.com/avito-tech-backend-autumn-2025/internal/delivery/http/handlers"
)

func NewRouter(
	teamHandler *handlers.TeamHandler,
	userHandler *handlers.UserHandler,
	prHandler *handlers.PRHandler,
	healthHandler *handlers.HealthHandler,
) *gin.Engine {
	r := gin.Default()

	teamHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)
	prHandler.RegisterRoutes(r)
	healthHandler.RegisterRoutes(r)

	return r
}
