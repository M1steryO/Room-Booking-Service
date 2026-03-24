include local.env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"


install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0 # для миграций
    GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@latest


local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

create-migration-example:
	${LOCAL_BIN}/goose create create_table_users sql -dir=./migrations



up:
	docker compose up --build -d

down:
	docker compose down -v

logs:
	docker compose logs -f app

# Наполнение базы тестовыми данными
seed:
	go run ./cmd/room-booking-service seed

test:
	go test ./... -coverprofile=coverage.out -covermode=atomic -coverpkg=./internal/...

test-e2e:
	@set -e; \
	trap 'docker compose -f docker-compose.e2e.yml down -v --remove-orphans' EXIT; \
	docker compose -f docker-compose.e2e.yml down -v --remove-orphans; \
	docker compose -f docker-compose.e2e.yml up -d postgres_e2e; \
	docker compose -f docker-compose.e2e.yml run migrator_e2e; \
	docker compose -f docker-compose.e2e.yml up -d --build app_e2e; \
	docker compose -f docker-compose.e2e.yml run tests


test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=./internal/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

lint:
	golangci-lint run ./...

swagger:
	swag init -g ./cmd/room-booking-service/main.go -o ./docs

run:
	go run ./cmd/room-booking-service serve

fmt:
	gofmt -w ./cmd ./internal ./tests



.PHONY: up down logs seed test test-e2e lint swagger run fmt
