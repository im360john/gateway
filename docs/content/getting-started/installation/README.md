---
title: 'Installation'
---

This guide provides three ways to install Gateway:

## 1. Binary Installation

Choose your operating system below for specific installation instructions for Linux:


```bash
# Download the latest binary for Linux
wget https://github.com/centralmind/gateway/releases/download/v0.1.1/gateway-linux-amd64.tar.gz

# Extract the archive
tar -xzf gateway-linux-amd64.tar.gz
mv gateway-linux-amd64 gateway

# Make the binary executable
chmod +x gateway

```

<details>
<summary>Windows (Intel)</summary>

```powershell
# Download the latest binary for Windows
Invoke-WebRequest -Uri https://github.com/centralmind/gateway/releases/download/v0.1.1/gateway-windows-amd64.zip -OutFile gateway-windows.zip

# Extract the archive
Expand-Archive -Path gateway-windows.zip -DestinationPath .

# Rename
Rename-Item -Path "gateway-windows-amd64.exe" -NewName "gateway.exe"

```
</details>

<details>
<summary>macOS (Intel)</summary>

```bash
# Download the latest binary for macOS (Intel)
curl -LO https://github.com/centralmind/gateway/releases/download/v0.1.1/gateway-darwin-amd64.tar.gz

# Extract the archive
tar -xzf gateway-darwin-amd64.tar.gz
mv gateway-darwin-amd64 gateway

# Make the binary executable
chmod +x gateway

```
</details>

<details>
<summary>macOS (Apple Silicon)</summary>
 
```bash
# Download the latest binary for macOS (Apple Silicon)
curl -LO https://github.com/centralmind/gateway/releases/download/v0.1.1/gateway-darwin-arm64.tar.gz

# Extract the archive
tar -xzf gateway-darwin-arm64.tar.gz
mv gateway-darwin-arm64 gateway

# Make the binary executable
chmod +x gateway

```

</details>


## 2. Docker Installation

### Option A: Use the pre-built Docker image

```bash
# Pull the latest image
docker pull ghcr.io/centralmind/gateway:v0.1.1

# Run the container
docker run -p 8080:8080 ghcr.io/centralmind/gateway:v0.1.1
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
./gateway --help
```
