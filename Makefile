.PHONY: build run dev test vet migrate seed clean docker-up docker-down

APP_NAME=rrhhumand
MAIN_PATH=./cmd/api

build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

run: build
	./bin/$(APP_NAME)

dev:
	go run $(MAIN_PATH)/main.go

test:
	go test ./... -v -count=1

test-coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

vet:
	go vet ./...

lint:
	golangci-lint run ./...

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-up:
	migrate -path migrations -database "postgres://postgres:Mfcd62!!Mfcd62!!@localhost:5432/rrhhumand?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:Mfcd62!!Mfcd62!!@localhost:5432/rrhhumand?sslmode=disable" down 1

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

clean:
	rm -rf bin/ coverage.out coverage.html
