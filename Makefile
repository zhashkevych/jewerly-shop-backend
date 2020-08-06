.PHONY: deploy

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/api/main.go

run: build
	docker-compose up --remove-orphans --build server

create-migration:
	migrate create -ext sql -dir schema/ -seq $(NAME)

deploy:
	export HOST=prod
	docker image build -t jewerly-api:0.1 -f ./deploy/Dockerfile .
	chmod +x ./deploy.sh
	./deploy.sh

migrate:
	migrate -path ./schema -database postgres://postgres:@0.0.0.0:55432/jewelryshop?sslmode=disable up

migrate-down:
	migrate -path ./schema -database postgres://postgres:@0.0.0.0:55432/jewelryshop?sslmode=disable down 1

migrate-drop:
	migrate -path ./schema -database postgres://postgres:@0.0.0.0:55432/jewelryshop?sslmode=disable drop
