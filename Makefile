run:
	@air

dev:
	@npm run dev

deps:
	@go mod tidy
	@go install github.com/a-h/templ/cmd/templ@latest

migrate-up:
	@go run ./cmd/migrations up

migrate-down:
	@go run ./cmd/migrations down

migration:
	@go run ./cmd/migrations create $(name)

asset:
	@npm run build
