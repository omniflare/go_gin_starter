FROM golang:1.23.4-alpine3.21

WORKDIR /src/app

RUN go install github.com/air-verse/air@latest

COPY . .

RUN go mod tidy
