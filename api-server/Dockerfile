FROM golang:1.18-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .

RUN go mod download
RUN go build -o api-server ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-server .
COPY --from=builder /app/sqlite/data /data

EXPOSE 8080

CMD ["./api-server"]