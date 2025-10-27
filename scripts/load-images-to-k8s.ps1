param(
    [string]$Tag = "dev"
)

$ErrorActionPreference = "Stop"

$services = @("auth", "game", "gateway")
$images = $services | ForEach-Object { "greatestworks-${_}:${Tag}" }

# Add infra images
$images += @("mongo:7", "redis:7")

Write-Host "Saving images to tar archives..." -ForegroundColor Cyan
foreach ($img in $images) {
    $safeName = $img -replace ":", "_"
    $tarFile = "tmp-${safeName}.tar"
    
    Write-Host "  Saving $img -> $tarFile" -ForegroundColor Gray
    docker save -o $tarFile $img
}

Write-Host "`nLoading images into kubernetes nodes..." -ForegroundColor Cyan
$nodes = kubectl get nodes -o jsonpath='{.items[*].metadata.name}'
foreach ($node in $nodes.Split(" ")) {
    Write-Host "  Node: $node" -ForegroundColor Yellow
    
    foreach ($img in $images) {
        $safeName = $img -replace ":", "_"
        $tarFile = "tmp-${safeName}.tar"
        
        Write-Host "    Loading $img" -ForegroundColor Gray
        docker cp $tarFile "${node}:/var/lib/${tarFile}"
        docker exec $node ctr -n k8s.io images import "/var/lib/${tarFile}"
        docker exec $node rm "/var/lib/${tarFile}"
    }
}

Write-Host "`nCleaning up local tar files..." -ForegroundColor Cyan
foreach ($img in $images) {
    $safeName = $img -replace ":", "_"
    $tarFile = "tmp-${safeName}.tar"
    Remove-Item $tarFile -ErrorAction SilentlyContinue
}

Write-Host "Done! Images are now available in the kubernetes cluster." -ForegroundColor Green
Write-Host "You can now deploy with imagePullPolicy: Never or IfNotPresent" -ForegroundColor Green
