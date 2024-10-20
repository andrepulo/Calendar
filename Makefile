BINARY_NAME=calendar
GO=go

.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME) ./cmd/main.go

.PHONY: run
run:
	$(GO) run ./cmd/main.go

.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down

.PHONY: docker-build
docker-build:
	docker build -t $(BINARY_NAME) .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)