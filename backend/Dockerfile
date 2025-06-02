# Use the official Go 1.24.3 image as the base image
FROM golang:1.24.3-alpine AS builder

# Set the working directory in the container
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o backend .

# Create a lightweight production image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/backend /app/backend

# Copy configuration files
COPY --from=builder /app/conf /app/conf
COPY --from=builder /app/static /app/static
COPY --from=builder /app/views /app/views

# Expose the port the app runs on
EXPOSE 8080

# Set environment variables
ENV RUNMODE=prod

# Command to run the application
CMD ["/app/backend"]
