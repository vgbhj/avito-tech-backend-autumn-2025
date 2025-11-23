package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/vgbhj/avito-tech-backend-autumn-2025/docs"

	"github.com/avito-tech-backend-autumn-2025/internal/delivery/http/handlers"
)

func NewRouter(
	teamHandler *handlers.TeamHandler,
	userHandler *handlers.UserHandler,
	prHandler *handlers.PRHandler,
	healthHandler *handlers.HealthHandler,
) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	teamHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)
	prHandler.RegisterRoutes(r)
	healthHandler.RegisterRoutes(r)

	return r
}
