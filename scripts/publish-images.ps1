param(
    [Parameter(Mandatory = $true)] [string]$Registry,          # e.g. docker.io
    [Parameter(Mandatory = $true)] [string]$Namespace,         # e.g. your-dockerhub-username
    [string]$Tag = "dev",                                    # image tag to publish
    [switch]$IncludeInfra,                                     # also re-tag/push mongo/redis
    [string]$MongoImage = "mongo:7",                          # infra image (optional)
    [string]$RedisImage = "redis:7"                           # infra image (optional)
)

$ErrorActionPreference = "Stop"

function Ensure-LoggedIn() {
    try {
        docker info | Out-Null
    } catch {
        Write-Host "Docker doesn't seem to be running. Please start Docker Desktop and try again." -ForegroundColor Yellow
        throw
    }
}

function Publish-ServiceImage([string]$service) {
    $local = "greatestworks-${service}:${Tag}"
    $remote = "${Registry}/${Namespace}/greatestworks-${service}:${Tag}"

    Write-Host "Pushing $local -> $remote" -ForegroundColor Cyan

    # Ensure local exists
    docker image inspect $local | Out-Null

    # Re-tag and push
    docker tag $local $remote
    docker push $remote
}

function Publish-InfraImage([string]$image) {
    # $image like: mongo:7 or redis:7
    $parts = $image.Split(":")
    $name = $parts[0]
    $tag = if ($parts.Length -gt 1) { $parts[1] } else { "latest" }

    $remote = "${Registry}/${Namespace}/${name}:${tag}"
    Write-Host "Re-tagging infra $image -> $remote" -ForegroundColor Cyan

    # Pull if missing locally (pull may fail if your Docker daemon has an invalid mirror)
    try {
        docker image inspect $image | Out-Null
    } catch {
        Write-Host "Local image $image not found; attempting to pull..." -ForegroundColor Yellow
        docker pull $image
    }

    docker tag $image $remote
    docker push $remote
}

Ensure-LoggedIn

$services = @("auth","game","gateway")
foreach ($svc in $services) { Publish-ServiceImage $svc }

if ($IncludeInfra) {
    Publish-InfraImage $MongoImage
    Publish-InfraImage $RedisImage
}

Write-Host "Done. Use the kustomize overlay at k8s/local/overlays/registry to deploy with these image names." -ForegroundColor Green
