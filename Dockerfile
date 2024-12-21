# Step 1: Build stage
FROM golang:alpine3.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies (using Go modules)
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Run stage
FROM alpine:latest

# Install certificates (if your Go app requires HTTPS)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=builder /app/main .

# Expose the port your app runs on (change as needed)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
