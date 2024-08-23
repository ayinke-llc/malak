tidy_dependencies:
	go mod tidy

migrate:
	migrate -path ./internal/datastore/postgres/migrations/ -database "postgres://malak:malak@localhost:9432/malak?sslmode=disable" up

migrate-down:
	migrate -path ./internal/datastore/postgres/migrations/ -database "postgres://malak:malak@localhost:9432/malak?sslmode=disable" down

run:
	go run cmd/*.go http

env:
	infisical export --env=dev > .env

test:
	go test -v ./...

test-all:
	go test -v -tags integration ./... 
