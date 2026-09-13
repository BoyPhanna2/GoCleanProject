# Build stage
FROM golang:1.27.1-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and certificates
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source code
COPY . .

# Build the application.
# Disable CGO to ensure a statically linked binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp ./cmd/api/main.go

# Final stage
FROM alpine:latest  

# Add tzdata for time-related operations (like logging)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/myapp .

# Copy swagger file
COPY --from=builder /app/swagger.yaml .

# Expose port 8080 to the outside world
EXPOSE 8080

# Set environment variables with defaults (can be overridden at runtime)
ENV PORT=8080
ENV DB_PATH=/data/myapp.db
ENV LOG_DIR=/logs

# Create volumes for persistent data
VOLUME ["/data", "/logs"]

# Command to run the executable
CMD ["./myapp"]
