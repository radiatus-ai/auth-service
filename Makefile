.PHONY: build run test clean

build-docker:
	docker build -t auth-service .

build:
	go build -o bin/server cmd/server/main.go

start:
	go run cmd/server/main.go

start-docker:
	docker compose --env-file .env.docker-local up api database --build

test:
	go test ./...

clean:
	rm -rf bin/

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

