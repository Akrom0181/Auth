migrate-create:
	migrate create -ext sql -dir migrations -seq auth

migrate-up:
	migrate -database 'postgres://akromjonotaboyev:1@localhost:5432/auth?sslmode=disable' -path migrations up

migrate-down:
	migrate -database 'postgres://akromjonotaboyev:1@localhost:5432/auth?sslmode=disable' -path migrations down

swag_init:
	swag init -g api/routers.go -o api/docs

run: swag_init
	go run cmd/main.go	

dev:
	nodemon --signal SIGTERM --watch . --ext go --exec "go run cmd/main.go"


	