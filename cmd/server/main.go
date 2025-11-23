package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avito-tech-backend-autumn-2025/api"
	"github.com/avito-tech-backend-autumn-2025/internal/config"
	"github.com/avito-tech-backend-autumn-2025/internal/database"
	"github.com/avito-tech-backend-autumn-2025/internal/delivery/http/handlers"
	"github.com/avito-tech-backend-autumn-2025/internal/domain"
	"github.com/avito-tech-backend-autumn-2025/internal/repository/postgres"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/pr"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/team"
	"github.com/avito-tech-backend-autumn-2025/internal/usecase/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	teamRepo := postgres.NewTeamRepository(db.DB)
	userRepo := postgres.NewUserRepository(db.DB)
	prRepo := postgres.NewPRRepository(db.DB)

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

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %d", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
