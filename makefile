# Makefile for create, up and down migrations in local and docker compose
# DATABASE_URL is the connection string like: postgres://postgres:password@localhost:5432/creditdb?sslmode=disable
# NAME is a simple string without spaces and sequence like: init_schema


## LOCAL

# migrations
migrations-up:
	migrate -path ./migrations -database $(DATABASE_URL) up

# migrations-down
migrations-down:
	migrate -path ./migrations -database $(DATABASE_URL) down

# create-migration
create-migration:
	migrate create -ext sql -dir ./migrations -seq $(NAME)

# get the current version of the migrations
migrations-version:
	migrate -path ./migrations -database $(DATABASE_URL) version

# force a specific version of the migrations
migrations-force:
	migrate -path ./migrations -database $(DATABASE_URL) force $(VERSION)


## DOCKER COMPOSE

# Start docker compose
compose-up:
	docker compose -f ./docker-compose.yml up -d --build

# Stop docker compose
compose-down:
	docker compose -f ./docker-compose.yml down -v

# Reset docker compose
compose-reset:
	docker compose -f ./docker-compose.yml down -v
	docker compose -f ./docker-compose.yml up -d --build


## REDIS LOCAL

# Start redis
redis-up:
	docker compose -f ./docker-compose.redis.local.yml up -d

# Stop redis
redis-down:
	docker compose -f ./docker-compose.redis.local.yml down -v

# Reset redis
redis-reset:
	docker compose -f ./docker-compose.redis.local.yml down -v
	docker compose -f ./docker-compose.redis.local.yml up -d


## GO LINT

# Run golangci-lint
go-lint:
	golangci-lint run ./...

# Run golangci-lint with fix
go-lint-fix:
	golangci-lint run ./... --fix


## GO TESTS

# Run unit tests
go-test:
	go test ./... -v -count=1 -short


## GO BENCHMARKS

# Run benchmarks
go-bench:
	go test ./internal/service/... -bench=. -benchmem -run=^$ -count=1


## GO BUILD AND RUN

# Build the binary
go-build:
	go build -o server ./cmd/server

# Run the binary
go-run:
	go run ./cmd/server

# Build and run the binary
go-build-run: go-build go-run
