FROM golang:1.14.0-alpine3.11 AS builder

RUN go version
RUN apk add git

COPY ./ /github.com/zhashkevych/jewelry-shop-backend
WORKDIR /github.com/zhashkevych/jewelry-shop-backend

RUN go mod download && go get -u ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/api/main.go

#lightweight docker container with binary
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=0 /github.com/zhashkevych/jewelry-shop-backend/.bin/app .
COPY --from=0 /github.com/zhashkevych/jewelry-shop-backend/pkg/config/ ./config/
COPY --from=0 /github.com/zhashkevych/jewelry-shop-backend/templates/ ./templates/

EXPOSE 8000

CMD [ "./app"]