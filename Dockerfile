# Use the official Go image as the base image
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a smaller base image for the final container
FROM debian:12-slim

# Set the Current Working Directory inside the container
WORKDIR /

# Copy the pre-built binary file from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/index.tmpl .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
