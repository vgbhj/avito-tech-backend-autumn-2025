package helpers

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func SetupTestDB() (*sql.DB, func(), error) {
	dbHost := getEnv("TEST_DB_HOST", "localhost")
	dbPort := getEnvAsInt("TEST_DB_PORT", 5432)
	dbUser := getEnv("TEST_DB_USER", "postgres")
	dbPassword := getEnv("TEST_DB_PASSWORD", "postgres")
	dbName := getEnv("TEST_DB_NAME", "pr_reviewer_test")

	adminDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)

	adminDB, err := sql.Open("postgres", adminDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to admin DB: %w", err)
	}
	defer adminDB.Close()

	_, _ = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))

	_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create test DB: %w", err)
	}

	testDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	testDB, err := sql.Open("postgres", testDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to test DB: %w", err)
	}

	if err := applyMigrations(testDB); err != nil {
		testDB.Close()
		return nil, nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	cleanup := func() {
		testDB.Close()
		adminDB, err := sql.Open("postgres", adminDSN)
		if err == nil {
			defer adminDB.Close()
			_, _ = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
		}
	}

	return testDB, cleanup, nil
}

func applyMigrations(db *sql.DB) error {
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS teams (
		team_name VARCHAR(255) PRIMARY KEY,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS users (
		user_id VARCHAR(255) PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		team_name VARCHAR(255) NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
		is_active BOOLEAN NOT NULL DEFAULT true,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name);
	CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
	CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active);

	CREATE TABLE IF NOT EXISTS pull_requests (
		pull_request_id VARCHAR(255) PRIMARY KEY,
		pull_request_name VARCHAR(255) NOT NULL,
		author_id VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
		status VARCHAR(50) NOT NULL DEFAULT 'OPEN',
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		merged_at TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_pr_author_id ON pull_requests(author_id);
	CREATE INDEX IF NOT EXISTS idx_pr_status ON pull_requests(status);
	CREATE INDEX IF NOT EXISTS idx_pr_author_status ON pull_requests(author_id, status);

	CREATE TABLE IF NOT EXISTS pr_reviewers (
		pull_request_id VARCHAR(255) NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
		reviewer_id VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
		assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
		PRIMARY KEY (pull_request_id, reviewer_id)
	);

	CREATE INDEX IF NOT EXISTS idx_pr_reviewers_reviewer_id ON pr_reviewers(reviewer_id);
	CREATE INDEX IF NOT EXISTS idx_pr_reviewers_pr_id ON pr_reviewers(pull_request_id);
	`

	_, err := db.Exec(migrationSQL)
	return err
}

func CleanupDB(db *sql.DB) error {
	_, err := db.Exec(`
		TRUNCATE TABLE pr_reviewers CASCADE;
		TRUNCATE TABLE pull_requests CASCADE;
		TRUNCATE TABLE users CASCADE;
		TRUNCATE TABLE teams CASCADE;
	`)
	return err
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
