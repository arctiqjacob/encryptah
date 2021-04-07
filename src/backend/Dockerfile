# Use the offical golang image to create a binary.
FROM golang:1.15-buster as builder

WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY app .
RUN go mod download

# Build the binary.
RUN go build -v -o server

# Build a small image
FROM debian:buster-slim

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

# Listens on port 5678
EXPOSE 5678

# Run the web service on container startup.
CMD ["/app/server"]