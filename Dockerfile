# Stage 1: Build Go binary
FROM golang:1.22-alpine AS builder

# Set up environment variables for Go
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o /gw .

# Stage 2: Base image setup (use Ubuntu for the other tools and dependencies)
FROM ubuntu:jammy

ENV TZ=Etc/UTC, ROTATION_TZ=Etc/UTC

ENV DEBIAN_FRONTEND=noninteractive

RUN echo $TZ > /etc/timezone && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime

# Create a non-root user and group
RUN addgroup --system cligroup && adduser --system --ingroup cligroup cliuser

# Copy the Go binary from Stage 1 (builder)
COPY --from=builder /gw /usr/local/bin/gw

RUN chmod +x /usr/local/bin/gw

# Set ownership of the binary to the non-root user
RUN chown cliuser:cligroup /usr/local/bin/gw

# Switch to the non-root user
USER cliuser

ENTRYPOINT ["/usr/local/bin/gw"]
