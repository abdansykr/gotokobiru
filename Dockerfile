# Stage 1: Build the Go application and the seeder using Go 1.23
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files first
COPY go.mod go.sum ./

# Copy the vendor directory which contains all dependencies.
COPY vendor ./vendor

# Copy the rest of the source code
COPY . .

# Build the main application binary using the vendored modules.
# The -mod=vendor flag is crucial here.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o /app/main ./main.go

# Build the seeder binary using the vendored modules.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o /app/seeder ./seed/seeder.go


# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /root/

# Copy the built binaries from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/seeder .

# Copy environment file template
COPY .env.example .

EXPOSE 8080

CMD ["./main"]
