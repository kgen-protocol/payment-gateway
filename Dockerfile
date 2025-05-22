# Stage 1: Build
FROM golang:1.23.4-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./cmd

# Stage 2: Deployment
FROM alpine:latest

# Set work directory and copy binary from builder
WORKDIR /root/
COPY --from=builder /app .

# Expose the application port
EXPOSE 8081
