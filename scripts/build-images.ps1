Param(
  [string]$Tag = "dev"
)

$ErrorActionPreference = 'Stop'

Write-Host "Building Docker images for greatestworks (tag=$Tag)"

# Resolve repo root
$RepoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path | Split-Path -Parent
Set-Location $RepoRoot

function Build-ServiceImage {
  param(
    [Parameter(Mandatory=$true)][string]$ServiceName,
    [Parameter(Mandatory=$true)][string]$ServicePackage
  )

  $image = "greatestworks-$ServiceName:$Tag"
  Write-Host " -> Building $image from package $ServicePackage"
  docker build `
    --build-arg SERVICE_PACKAGE=$ServicePackage `
    --build-arg BUILD_VERSION=$Tag `
    --build-arg BUILD_TIME=$(Get-Date -Format o) `
    --build-arg GIT_COMMIT=$(git rev-parse --short HEAD 2>$null) `
    -t $image `
    -f Dockerfile .
}

Build-ServiceImage -ServiceName "auth" -ServicePackage "./cmd/auth-service"
Build-ServiceImage -ServiceName "game" -ServicePackage "./cmd/game-service"
Build-ServiceImage -ServiceName "gateway" -ServicePackage "./cmd/gateway-service"

Write-Host "Done. Images:"
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | Select-String greatestworks-
