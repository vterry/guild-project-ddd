build:
	@go build -o bin/character cmd/main.go

migration:
	@migrate create -ext sql -dir internal/infra/db/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run internal/infra/db/migrate/main.go up

migrate-down:
	@go run internal/infra/db/migrate/main.go down

migrate-reset:
	@go run internal/infra/db/migrate/main.go force

test:
	@go test -v ./...

test-cover:
	@go test -count=1 -coverprofile coverage.out ./... && go tool cover -func=coverage.out

run: build
	@./bin/character
