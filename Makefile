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

# canvas/ada cli will be able to replace this file
build-docker:
	docker compose build api-deploy

tag:
	docker tag auth-service-api-deploy:latest us-central1-docker.pkg.dev/rad-containers-hmed/shared/auth-service:latest

upload:
	docker push us-central1-docker.pkg.dev/rad-containers-hmed/shared/auth-service:latest

# todo: move to k8s cluster
deploy: build-docker tag upload
	gcloud run deploy auth-service \
					--image=us-central1-docker.pkg.dev/rad-containers-hmed/shared/auth-service:latest \
					--execution-environment=gen2 \
					--region=us-central1 \
					--project=rad-dev-platapi-4r64 \
					&& gcloud run services update-traffic auth-service --to-latest --region us-central1 --project=rad-dev-platapi-4r64
