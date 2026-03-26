# Use official Go image as the build environment
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Use a minimal image for running
FROM alpine:latest
# Set the working directory
WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Command to run
CMD ["./main"]