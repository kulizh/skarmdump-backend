APP=skarmdump-backend
IMAGE=skarmdump-backend

build:
	go build -o $(APP) ./cmd/server

build-ubuntu:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(APP) ./cmd/server

run:
	go run ./cmd/server

# Локальный запуск через Docker
docker-build:
	docker build -t $(IMAGE) .

docker-run: docker-build
	docker run --rm -p 8080:8080 \
		--env-file .env \
		-v ./img:/app/img \
		$(IMAGE)

compose-up:
	docker compose up -d --build

compose-down:
	docker compose down

compose-logs:
	docker compose logs -f