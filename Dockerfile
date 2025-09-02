# Use Go 1.22 as the base image
FROM golang:1.22

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Install swag CLI for Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/main.go  

# Build the application
RUN go build -o bin/api cmd/main.go

# Expose the port (Railway assigns PORT, defaults to 8080)
EXPOSE 8080

# Set environment variables (overridden by Railway)
ENV PORT=8080
ENV BASE_URL=http://localhost:8080

# Run the application
CMD ["./bin/api"]