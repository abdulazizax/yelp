run :
	go run cmd/main.go

build:
	docker-compose up --build -d

down:
	docker-compose down

swag_init:
	swag init -g internal/adapter/http/router.go --parseDependency -o api/docs