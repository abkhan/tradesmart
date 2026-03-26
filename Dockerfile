# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install git for dependencies
RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./
# Run tidy to ensure all dependencies are resolved in the build context
RUN go mod tidy
RUN go mod download

# Copy the rest of the source code
COPY . .

# Argument to decide which app to build (orders-api or reports-api)
ARG APP_NAME
RUN go build -o /app/service ./cmd/${APP_NAME}/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/service .

# Standard ports
EXPOSE 8080 8081

# Run the service
CMD ["./service"]
