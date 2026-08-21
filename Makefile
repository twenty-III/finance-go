build:
	@go build -o bin/finance cmd/api/main.go

test:
	@go test -v ./...

run: build
	@./bin/finance

gen:
	@go run github.com/99designs/gqlgen generate
	
