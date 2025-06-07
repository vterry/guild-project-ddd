build:
	@go build -o bin/character cmd/main.go

test:
	@go test -v ./...

test-cover:
	@go test -count=1 -coverprofile coverage.out ./... && go tool cover -func=coverage.out

run: build
	@./bin/character
