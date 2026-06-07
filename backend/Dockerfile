# Stage 1: Build
FROM golang:alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Stage 2: Run
FROM alpine:latest
ENV TZ=Asia/Jakarta
RUN apk add --no-cache tzdata

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/main .
# Copy migrations and API docs
COPY --from=builder /app/database/migrations ./migrations
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./main"]
