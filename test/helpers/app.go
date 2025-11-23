package helpers

import (
	"database/sql"

	"github.com/avito-tech-backend-autumn-2025/api"
	"github.com/avito-tech-backend-autumn-2025/internal/delivery/http/handlers"
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/postgres"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/pr"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/team"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

func SetupTestApp(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	teamRepo := postgres.NewTeamRepository(db)
	userRepo := postgres.NewUserRepository(db)
	prRepo := postgres.NewPRRepository(db)

	reviewerAssigner := domain.NewReviewerAssigner()

	createTeamUseCase := team.NewCreateTeamUseCase(teamRepo, userRepo)
	getTeamUseCase := team.NewGetTeamUseCase(teamRepo)
	setActiveUseCase := user.NewSetActiveUseCase(userRepo)
	getReviewsUseCase := user.NewGetReviewsUseCase(prRepo, userRepo)
	createPRUseCase := pr.NewCreatePRUseCase(prRepo, userRepo, teamRepo, reviewerAssigner)
	mergePRUseCase := pr.NewMergePRUseCase(prRepo)
	reassignReviewerUseCase := pr.NewReassignReviewerUseCase(prRepo, userRepo, teamRepo, reviewerAssigner)

	teamHandler := handlers.NewTeamHandler(createTeamUseCase, getTeamUseCase)
	userHandler := handlers.NewUserHandler(setActiveUseCase, getReviewsUseCase)
	prHandler := handlers.NewPRHandler(createPRUseCase, mergePRUseCase, reassignReviewerUseCase)
	healthHandler := handlers.NewHealthHandler()

	router := api.NewRouter(teamHandler, userHandler, prHandler, healthHandler)

	return router
}
