# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /users

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o users_app ./cmd/main.go

FROM alpine:latest

WORKDIR /users

COPY --from=builder /users/users_app .

COPY . .

EXPOSE 8443 9443

CMD ["./users_app"]