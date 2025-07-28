# for my golang http/net server

FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .


# Final stage
FROM alpine:latest
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /app/server .
