.PHONY: docs test run build

docs:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/main.go
	go run ./cmd/openapi -input docs/swagger.json -output docs/openapi.json

test:
	go test ./...

run: docs
	SWAG_AUTO=0 go run ./cmd

build: docs
	mkdir -p bin
	go build -o bin/tower-go ./cmd
