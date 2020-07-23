run-db:
	docker run -d --name jewelry-shop-2 -v ./.build/data/:/var/lib/postgresql/data -p 54320:5432 postgres:12.0-alpine

run:
	go run cmd/api/main.go
