# Use official Go image as the build environment
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Use a minimal image for running
FROM alpine:latest

# Install sqlite for debugging (optional but helpful)
RUN apk add --no-cache sqlite

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/main .

# IMPORTANT: Create the data directory and set permissions
RUN mkdir -p /app/data && chmod 777 /app/data

# Expose the port
EXPOSE 8080

# Command to run
CMD ["./main"]