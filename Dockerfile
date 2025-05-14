# First stage: Build the application
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY app/go.mod app/go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY app/ ./

# Build the application
RUN go build -o iotgo-app .

# Second stage: Create the final image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install dependencies for PostgreSQL
RUN apk --no-cache add ca-certificates postgresql-client

# Copy the binary from the builder stage
COPY --from=builder /app/iotgo-app .

# Copy configuration files
COPY app/conf ./conf

# Copy static files and views
COPY app/static ./static
COPY app/views ./views

# Expose port (adjust if your app uses a different port)
EXPOSE 8080

# Set environment variables (these can be overridden at runtime)
ENV DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=iotnew \
    DB_HOST=postgres \
    DB_PORT=5432 \
    DB_SSLMODE=disable

# Run the application
CMD ["./iotgo-app"]