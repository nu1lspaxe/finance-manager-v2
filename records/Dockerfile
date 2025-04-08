# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /records

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o records_app ./cmd/main.go

FROM alpine:latest

WORKDIR /records

COPY --from=builder /records/records_app .

COPY . .

EXPOSE 8444 9444

CMD ["./records_app"]