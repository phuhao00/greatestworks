Param(
  [switch]$UseMinikube = $false,
  [string]$Namespace = "gaming",
  [string]$Tag = "dev"
)

$ErrorActionPreference = 'Stop'

# Repo root and kubectl presence
$RepoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path | Split-Path -Parent
Set-Location $RepoRoot

if (-not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
  throw "kubectl is required but not found in PATH"
}

# Ensure namespace and base infra
Write-Host "Applying namespace and base infra (MongoDB, Redis)"
kubectl apply -f k8s/local/namespace.yaml | Out-Host
kubectl apply -n $Namespace -f k8s/local/mongodb.yaml | Out-Host
kubectl apply -n $Namespace -f k8s/local/redis.yaml | Out-Host

# Optionally load local images into minikube
if ($UseMinikube) {
  if (-not (Get-Command minikube -ErrorAction SilentlyContinue)) {
    throw "--UseMinikube specified but 'minikube' command not found"
  }
  Write-Host "Loading local images into minikube: greatestworks-auth:$Tag, greatestworks-game:$Tag, greatestworks-gateway:$Tag"
  minikube image load "greatestworks-auth:$Tag"
  minikube image load "greatestworks-game:$Tag"
  minikube image load "greatestworks-gateway:$Tag"
}

# Apply ConfigMaps and services
Write-Host "Applying ConfigMaps and service deployments"
kubectl apply -n $Namespace -f k8s/local/configmap-gateway.yaml | Out-Host
kubectl apply -n $Namespace -f k8s/local/game-service.yaml | Out-Host
kubectl apply -n $Namespace -f k8s/local/auth-service.yaml | Out-Host
kubectl apply -n $Namespace -f k8s/local/gateway-service.yaml | Out-Host

# Wait for pods readiness
Write-Host "Waiting for deployments to be ready..."
kubectl -n $Namespace rollout status deployment/mongodb --timeout=120s | Out-Host
kubectl -n $Namespace rollout status deployment/redis --timeout=120s | Out-Host
kubectl -n $Namespace rollout status deployment/game-service --timeout=180s | Out-Host
kubectl -n $Namespace rollout status deployment/auth-service --timeout=180s | Out-Host
kubectl -n $Namespace rollout status deployment/gateway-service --timeout=180s | Out-Host

Write-Host "Services (NodePort)"
kubectl -n $Namespace get svc | Out-Host

Write-Host "Done. You can reach:"
Write-Host " - Auth HTTP:    http://localhost:30080"
Write-Host " - Gateway TCP:  <your host IP>:30909"
