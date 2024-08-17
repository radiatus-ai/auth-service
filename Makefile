.PHONY: build run test clean

build:
	go build -o bin/server cmd/server/main.go

start:
	go run cmd/server/main.go

test:
	go test ./...

clean:
	rm -rf bin/

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

docker-build:
	docker build -t auth-service .

docker-run:
	docker run -p 8080:8080 auth-service
