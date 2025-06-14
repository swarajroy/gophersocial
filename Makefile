include .envrc
MIGRATIONS_PATH=./cmd/migrate/migrations
DB_MIGRATOR_ADDR="postgres://admin:adminpassword@localhost/social?sslmode=disable"

.PHONY: test
test:
	@go test -v ./...

.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_MIGRATOR_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_MIGRATOR_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: seed
seed:
	@go run ./cmd/migrate/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g main.go -d cmd/api,internal/store && swag fmt