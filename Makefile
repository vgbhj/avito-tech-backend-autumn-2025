.PHONY: migrate test test-integration

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pr_reviewer_db?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pr_reviewer_db?sslmode=disable" down

test:
	go test ./...

test-integration:
	go test -v ./test/integration/...