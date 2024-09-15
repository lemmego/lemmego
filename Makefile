run:
	@air

dev:
	@npm run dev

build:
	@npm run build

deps:
	@go mod tidy
	@go install github.com/a-h/templ/cmd/templ@latest

migrate:
	@go run ./cmd/migrations up

rollback:
	@go run ./cmd/migrations down

migration:
	@go run ./cmd/migrations create $(n)

handlers:
	@go run ./cmd g handlers

input:
	@go run ./cmd g input

model:
	@go run ./cmd g model

form:
	@go run ./cmd g form
