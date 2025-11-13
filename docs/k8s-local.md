# Run GreatestWorks locally on Kubernetes (Windows)

This guide helps you run the Go services locally on a Kubernetes cluster using Docker Desktop or minikube on Windows PowerShell.

## Prerequisites

- Windows with PowerShell 5.1+
- Docker Desktop (recommended) or minikube
- kubectl in PATH
- (If using minikube) minikube in PATH

## Build images

Build three service images (auth, game, gateway) using the unified Dockerfile:

```
# From repo root
./scripts/build-images.ps1 -Tag dev
```

This produces:
- greatestworks-auth:dev
- greatestworks-game:dev
- greatestworks-gateway:dev

## Deploy to Kubernetes

Docker Desktop (Kubernetes enabled):

```
./scripts/k8s-deploy.ps1 -Namespace gaming -Tag dev
```

Minikube:

```
# Ensure 'minikube start' is already running
./scripts/k8s-deploy.ps1 -UseMinikube -Namespace gaming -Tag dev
```

The script applies:
- Namespace: `gaming`
- Infra: MongoDB (user: admin / pass: admin123), Redis (password: redis123)
- Services: auth-service (HTTP 8080), game-service (RPC 8081), gateway-service (TCP 9090)
- NodePorts: auth-service 30080, gateway-service 30909

## Verify

```
kubectl -n gaming get pods
kubectl -n gaming get svc
```

- Auth HTTP: http://localhost:30080/health
- Gateway TCP: connect to port 30909 from your client

## Notes

- Config files are bundled in the image under `/configs`.
- Auth-service DB URIs and Redis address are overridden via env vars in the Deployment.
- Gateway-service config is provided via a ConfigMap mounted to `/configs/gateway-service.yaml` to point at in-cluster services.
- For non-NodePort clusters, use `kubectl port-forward` as needed.

## Cleanup

```
kubectl delete ns gaming
```
