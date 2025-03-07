# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /bank_system

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bank_system ./cmd/main.go

FROM alpine:latest

WORKDIR /bank_system

COPY --from=builder /bank_system/bank_system .

COPY configs/config.json ./configs/config.json
COPY certs/cert.pem ./certs/cert.pem
COPY certs/key.pem ./certs/key.pem

EXPOSE 8855

CMD ["./bank_system"]