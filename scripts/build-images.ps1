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

  $imageName = "greatestworks-${ServiceName}:${Tag}"
  Write-Host " -> Building $imageName from package $ServicePackage"
  $buildTime = Get-Date -Format o
  $gitCommit = git rev-parse --short HEAD 2>$null
  docker build `
    --build-arg SERVICE_PACKAGE=$ServicePackage `
    --build-arg BUILD_VERSION=$Tag `
    --build-arg BUILD_TIME=$buildTime `
    --build-arg GIT_COMMIT=$gitCommit `
    -t $imageName `
    -f Dockerfile .
}

Build-ServiceImage -ServiceName "auth" -ServicePackage "./cmd/auth-service"
Build-ServiceImage -ServiceName "game" -ServicePackage "./cmd/game-service"
Build-ServiceImage -ServiceName "gateway" -ServicePackage "./cmd/gateway-service"

Write-Host "Done. Images:"
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | Select-String greatestworks-
