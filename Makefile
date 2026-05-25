DATABASE_URL=postgres://admin:supersecretpassword@127.0.0.1:5432/blog?sslmode=disable

APP := bin/api
ENV := .env

run:
	@env $$(cat $(ENV) | grep -v '^#' | xargs) go run .

build:
	@mkdir -p bin && go build -o $(APP) .

test:
	go test ./...

test-race:
	go test -race ./...

tidy:
	go mod tidy

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	staticcheck ./...

vuln:
	govulncheck ./...

check:
	@test -z "$$(gofmt -l .)" || (echo "run make fmt first"; exit 1)
	@$(MAKE) vet test test-race lint vuln
	@echo "✓ all checks passed"

up:
	docker compose up -d

down:
	docker compose down

nuke:
	docker compose down -v

logs:
	docker compose logs -f

psql:
	docker exec -it blog-postgres-1 psql -U admin -d blog

migrate-new:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-version:
	migrate -path migrations -database "$(DATABASE_URL)" version
