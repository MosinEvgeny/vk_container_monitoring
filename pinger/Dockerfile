FROM golang:1.23.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env .
RUN go build -o pinger ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/pinger .
COPY .env .
CMD ["./pinger"]