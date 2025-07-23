# for my golang http/net server

FROM golang:1.24.4-alpine AS builder
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
# Download dependencies
RUN go mod download
# Copy the source code
COPY . .
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o server .
# Final stage
FROM alpine:latest
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /app/server .
# Command to run the binary
# CMD ["./server"]
