# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install git for module fetches
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

RUN chmod +x ./app

CMD ["./app"]
