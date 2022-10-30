.PHONY:
.SILENT:

lint:
	staticcheck ./...
	golangci-lint run

test:
	go test --short -coverprofile=cover.out -v ./...
	make test.coverage

test.coverage:
	go tool cover -func=cover.out | grep "total"

# go install github.com/swaggo/swag/cmd/swag@latest
swag:
	swag init -g cmd/main.go

run:
	docker-compose up -d --build
