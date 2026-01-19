# Stage 1: Build the application
FROM golang:1.25.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o crawler main.go

# Stage 2: Final lightweight image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/crawler .
CMD ["./crawler"]