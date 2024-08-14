tidy_dependencies:
	go mod tidy

migrate:
	migrate -path ./internal/datastore/postgres/migrations/ -database "postgres://malak:malak@localhost:3432/malak?sslmode=disable" up

migrate-down:
	migrate -path ./internal/datastore/postgres/migrations/ -database "postgres://malak:malak@localhost:3432/malak?sslmode=disable" down

run:
	go run cmd/*.go http

env:
	infisical export --env=dev > .env
