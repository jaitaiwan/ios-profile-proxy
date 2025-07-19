# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for module fetches
RUN apk add --no-cache git

COPY . .

RUN go mod download || true
RUN go build -o app .

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

RUN chmod +x ./app

CMD ["./app"]
