# Gateway Demo Deployment Guide

This guide will help you deploy the Gateway demo application with PostgreSQL in Kubernetes.

## Prerequisites

- Kubernetes cluster
- Helm v3
- `kubectl` configured to work with your cluster
- Access to GitHub Container Registry (ghcr.io)

## Setup GitHub Container Registry Secret

Before deploying, create a secret for pulling images from GitHub Container Registry:

```bash
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=YOUR_GITHUB_USERNAME \
  --docker-password=YOUR_GITHUB_PAT \
  --namespace=demo
```

Replace `YOUR_GITHUB_USERNAME` with your GitHub username and `YOUR_GITHUB_PAT` with your GitHub Personal Access Token.

## Deployment Steps

1. Deploy PostgreSQL database:
```bash
make install-postgres
```

2. Wait for PostgreSQL to be ready:
```bash
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=postgresql -n demo
```

3. Deploy the Gateway application:
```bash
make install-gateway
```

## Verification

1. Check if all pods are running:
```bash
kubectl get pods -n demo
```

2. Access the API:
The API will be available at: `http://demo-gw.centralmind.ai`

Example endpoints:
- GET `/gachi_teams` - List all teams
- GET `/gachi_personas` - List all personas

## Useful Commands

- Get PostgreSQL password:
```bash
make get-password
```

- Upgrade Gateway configuration:
```bash
make upgrade-gateway
```

- Upgrade PostgreSQL configuration:
```bash
make upgrade-postgres
```

## Cleanup

To remove the deployment:

1. Uninstall Gateway:
```bash
make uninstall-gateway
```

2. Uninstall PostgreSQL:
```bash
make uninstall-postgres
```

## Configuration

The deployment uses two main configuration files:
- `values.gateway.yaml` - Gateway configuration
- `values.postgres.yaml` - PostgreSQL configuration

Modify these files to customize your deployment. 