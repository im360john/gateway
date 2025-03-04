---
title: 'Installation'
---

This guide provides three ways to install Gateway:

## 1. Binary Installation

The simplest method is to download and use the pre-compiled binary:

```bash
# Download the latest binary
wget https://github.com/centralmind/gateway/releases/latest/download/gateway_linux_amd64.tar.gz

# Extract the archive
tar -xzf gateway_linux_amd64.tar.gz

# Make the binary executable
chmod +x gateway

# Move to a directory in your PATH (optional)
sudo mv gateway /usr/local/bin/
```

## 2. Docker Installation

### Option A: Use the pre-built Docker image

```bash
# Pull the latest image
docker pull ghcr.io/centralmind/gateway:latest

# Run the container
docker run -p 8080:8080 ghcr.io/centralmind/gateway:latest
```

### Option B: Build your own Docker image

```bash
# Clone the repository
git clone https://github.com/centralmind/gateway.git
cd gateway

# Build the Docker image
docker build -t gateway .

# Run your container
docker run -p 8080:8080 gateway
```

## 3. Manual Build from Source

To compile the application yourself:

```bash
# Clone the repository
git clone https://github.com/centralmind/gateway.git
cd gateway

# Install dependencies and build
go mod download
go build -o gateway

# Make the binary executable
chmod +x gateway

# Run the application
./gateway
```
