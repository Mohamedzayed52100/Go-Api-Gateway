# Use the latest Go version as your builder image
FROM golang:latest as builder

# Set the working directory inside the container
WORKDIR /build


# Your GitHub Access Token for private repositories (if needed)
ARG GITHUB_ACCESS_TOKEN

# Copy Go mod and sum files to leverage Docker layer caching
COPY go.mod go.sum ./

# Configure Git to use your access token
RUN git config --global url."https://${GITHUB_ACCESS_TOKEN}:@github.com/".insteadOf "https://github.com/"

# Download Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application with debug symbols (no optimizations)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd

# Run tests
# RUN go test -v ./...
# Use a lightweight base image for the final image
FROM alpine:latest


# Set the working directory in the container
WORKDIR /app

# Copy the built executable and Delve debugger from the builder image
COPY --from=builder /app/main ./

# Expose ports for the application and Delve debugger
# EXPOSE 8081 40000
EXPOSE 8081 

# Define command to run the application with Delve for debugging
CMD ["/app/main"]

