# Build stage
FROM golang:latest AS builder

RUN mkdir /build

WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest


# Copy the binary from builder
COPY --from=builder /build/main .

# Set the working environment variables




# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"] 