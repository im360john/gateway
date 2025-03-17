---
title: Deploy Gateway on Kubernetes
---

On this page, you will find instructions for installing and running the Gateway demo application in Kubernetes using Kubernetes manifests.

## Before you begin

To follow this guide:

* You need the latest version of Kubernetes running either locally or remotely on a public or private cloud.
* If you plan to use it in a local environment, you can use various Kubernetes options such as minikube, kind, Docker Desktop, and others.
* If you plan to use Kubernetes in a production setting, it's recommended to utilize managed cloud services like Google Kubernetes Engine (GKE), Amazon Elastic Kubernetes Service (EKS), or Azure Kubernetes Service (AKS).
* Configured `kubectl` to work with your cluster
* Access to GitHub Container Registry (ghcr.io)

## System requirements

This section provides minimum hardware and software requirements.

### Minimum Hardware Requirements

* Memory: 250 MiB (approx 250 MB)
* CPU: 125m (approx 0.125 cores)

> Note
> 
> Enable the necessary ports in your network environment for Gateway.

## Deploy Gateway on Kubernetes

This section explains how to install Gateway using Kubernetes.

### Create a namespace

It is recommended to create a new namespace in Kubernetes to better manage, organize, allocate, and manage cluster resources:

```bash
kubectl create namespace demo
```

### Deploy the Gateway application

1. Create a `gateway.yaml` file, and copy-and-paste the following content into it:
   ```yaml
   ---
   apiVersion: v1
   kind: ConfigMap
   metadata:
     name: gateway-config
     namespace: demo
   data:
     config.yaml: |
        api:
            name: Automatic API
            description: ""
            version: 1.0.0
        database:
             type: postgres
             connection: "connection-string"
        plugins: {}
   ---
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: gateway
     namespace: demo
     labels:
       app: gateway
   spec:
     replicas: 1
     selector:
       matchLabels:
         app: gateway
     template:
       metadata:
         labels:
           app: gateway
       spec:
         volumes:
           - name: config
             configMap:
               name: gateway-config
         containers:
         - name: gateway
           image: ghcr.io/centralmind/gateway:latest
           args:
               - start
               - --config
               - /etc/gateway/config.yaml
               - --addr
               - ":8080"
           imagePullPolicy: Always
           ports:
           - containerPort: 8080
             name: http
           volumeMounts:
            - name: config
              mountPath: /etc/gateway
              readOnly: true
   ---
   apiVersion: v1
   kind: Service
   metadata:
     name: gateway
     namespace: demo
   spec:
     ports:
     - port: 80
       targetPort: http
     selector:
       app: gateway
   ```

2. Run the following command to send the manifest to the Kubernetes API server:
   ```bash
   kubectl apply -f gateway.yaml
   ```

### Verify the deployment

1. Complete the following steps to verify the deployment status of each object:
   
   a. For Pods, run the following command:
   ```bash
   kubectl get pods -n demo
   ```
   
   b. For Services, run the following command:
   ```bash
   kubectl get svc -n demo -o wide
   ```

### Access the Gateway API

To access the Gateway API from your local machine, you can use port forwarding. This allows you to connect to the Gateway service running in your Kubernetes cluster without exposing it externally.

1. First, identify the Gateway pod name:
   ```bash
   kubectl get pods -n demo -l app=gateway
   ```

2. Set up port forwarding to the Gateway pod:
   ```bash
   kubectl port-forward -n demo svc/gateway 8080:8080
   ```
   This command forwards your local port 8080 to port 80 of the Gateway service in the Kubernetes cluster.

3. Now you can access the Gateway API at: `http://localhost:8080`

Example endpoints:
- GET `/gachi_teams` - List all teams
- GET `/gachi_personas` - List all personas

> Note
> 
> The port forwarding session will continue until you terminate it with Ctrl+C. If you close the terminal or the connection is interrupted, you'll need to restart the port forwarding.

## Configuration

You can modify the Kubernetes manifests to customize your deployment:
- `gateway.yaml` - Gateway configuration

## Update an existing deployment

### Update Gateway configuration

To update the Gateway configuration:

1. Edit the `gateway.yaml` file with your changes
2. Apply the updated configuration:
   ```bash
   kubectl apply -f gateway.yaml
   ```

## Troubleshooting

### Collecting logs

To collect Gateway logs, run:
```bash
kubectl logs -l app=gateway -n demo
```

### Using the --dry-run command

You can use the `--dry-run=client` flag with kubectl to test your manifests without actually applying them:

```bash
kubectl apply -f gateway.yaml --dry-run=client
```

## Remove Gateway

To remove the deployment:

1. Delete Gateway resources:
   ```bash
   kubectl delete -f gateway.yaml
   ```

2. Delete the namespace (optional):
   ```bash
   kubectl delete namespace demo
   ```

## Was this page helpful?

If you have questions or suggestions for improving this documentation, please create an issue in the project repository. 
