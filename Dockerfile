FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Step 2: Create the final, minimal image
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY migrations migrations

# Command to run the executable
CMD ["./main"]
