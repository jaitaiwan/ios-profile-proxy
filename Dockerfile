# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for module fetches
RUN apk add --no-cache git

COPY . .

RUN go mod download
RUN go build -o app .

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

RUN chmod +x ./app

CMD ["./app"]
