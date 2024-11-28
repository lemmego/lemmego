run:
	@air

templ-notify:
	@templ generate --notify-proxy

templ:
	@templ generate --watch --proxy="http://localhost:8080" -v

tailwind:
	npx --yes tailwindcss -i static/css/style.css -o static/css/dist.css --minify --watch

watch:
	make -j4 templ run templ-notify tailwind

dev:
	@npm run dev

build:
	@npm run build

deps:
	@go mod tidy
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install github.com/air-verse/air@latest

migrate:
	@lemmego run migrate up

rollback:
	@lemmego run migrate down

migration:
	@lemmego run migrate create $(n)

handlers:
	@lemmego g handlers $(n)

input:
	@lemmego g input $(n)

model:
	@lemmego g model $(n)

form:
	@lemmego g form $(n)