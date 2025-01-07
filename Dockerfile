# Stage 1: Build stage
FROM golang:1.23.4 as builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags migrate -o yelpbinary ./cmd/app

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary tools
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/yelpbinary /app/
# Copy configuration and migration files from the builder stage
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations

# Copy the .env file if it exists in the root directory
COPY .env .env

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/yelpbinary"]
