build:
	@go build -o bin/finance cmd/api/main.go

test:
	@go test -v ./...

run: build
	@./bin/finance

gen:
	@go run github.com/99designs/gqlgen generate

migrate-up:
	@goose -dir internal/db/migrations postgres "user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" up

migrate-down:
	@goose -dir internal/db/migrations postgres "user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" down
