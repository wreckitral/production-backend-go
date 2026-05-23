DATABASE_URL=postgres://admin:supersecretpassword@127.0.0.1:5432/blog?sslmode=disable

migrate-new:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1
