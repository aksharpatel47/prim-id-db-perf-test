install-deps:
	echo "Installing dependencies..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

create-migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -source file://migrations -database postgres://postgres:postgres@localhost:55433/postgres?sslmode=disable up

migrate-down:
	migrate -source file://migrations -database postgres://postgres:postgres@localhost:55433/postgres?sslmode=disable down