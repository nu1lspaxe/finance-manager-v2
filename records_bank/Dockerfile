# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /records_bank

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o records_bank_app ./cmd/main.go

FROM alpine:latest

WORKDIR /records_bank

COPY --from=builder /records_bank/records_bank_app .

COPY . .

EXPOSE 8445 9445

CMD ["./records_bank_app"]