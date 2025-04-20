build:
	@go build -o bin/auth-server cmd/main.go

test:
	@go test -v ./...

test-cover:
	@go test -count=1 -coverprofile coverage.out ./...

run: build
	@./bin/auth-server

