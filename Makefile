CURRENT_DIR=$(shell pwd)
DB_URL=postgres://postgres:pass@localhost:5432/twitter?sslmode=disable

proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}

run :
	go run cmd/main.go
  
migrate_up:
	migrate -path migrations -database ${DB_URL}  -verbose up

migrate_down:
	migrate -path migrations -database ${DB_URL}  -verbose down

migrate_force:
	migrate -path migrations -database ${DB_URL}  -verbose force 1

migrate_file:
	migrate create -ext sql -dir migrations -seq create_tables

test:
	go test -v -cover ./...

swag_init:
	swag init -g internal/http/app/app.go --parseDependency -o internal/http/app/docs