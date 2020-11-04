.SILENT:

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/api/main.go

run: build
	docker-compose up --remove-orphans --build server

test:
	go test ./... -coverprofile cover.out

test-coverage:
	go tool cover -func cover.out | grep total | awk '{print $3}'

create-migration:
	migrate create -ext sql -dir schema/ -seq $(NAME)

migrate:
	migrate -path ./schema -database postgres://postgres:@0.0.0.0:5432/postgres?sslmode=disable up

migrate-down:
	migrate -path ./schema -database postgres://postgres:@0.0.0.0:5432/postgres?sslmode=disable down 1

migrate-drop:
	migrate -path ./schema -database postgres://postgres:@0.0.0.0:5432/postgres?sslmode=disable drop

logs-stage:
	tail -f /root/jewerly-shop/api/logs/stage/api.log

logs-prod:
	tail -f /root/jewerly-shop/api/logs/prod/api.log