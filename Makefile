tidy_dependencies:
	go mod tidy

migrate:
	go run cmd/*.go migrate

migrate-down:
	migrate -path ./internal/datastore/postgres/migrations/ -database "postgres://malak:malak@localhost:9432/malak?sslmode=disable" down 1

run:
	go run cmd/*.go http

env:
	infisical export --env=dev > .env

test:
	go test -v ./...

test-all:
	go test -v -tags integration ./... 
