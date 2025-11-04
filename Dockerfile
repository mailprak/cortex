# Multi-stage build for Cortex
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cortex .

# Runtime image
FROM ubuntu:22.04

# Install kubectl and other common tools
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    bash \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl \
    && mv kubectl /usr/local/bin/

# Copy cortex binary from builder
COPY --from=builder /build/cortex /usr/local/bin/cortex

# Copy examples
COPY example /cortex/example

WORKDIR /cortex

# Set cortex as executable
RUN chmod +x /usr/local/bin/cortex

# Default command
CMD ["cortex", "--help"]
