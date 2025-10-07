# PowerShell脚本修复导入路径
Get-ChildItem -Path . -Filter "*.go" -Recurse | ForEach-Object {
    $content = Get-Content $_.FullName -Raw
    if ($content -match "greatestworks/internal/infrastructure/logger") {
        $newContent = $content -replace "greatestworks/internal/infrastructure/logger", "greatestworks/internal/infrastructure/logging"
        Set-Content -Path $_.FullName -Value $newContent -NoNewline
        Write-Host "Fixed: $($_.FullName)"
    }
}
Write-Host "导入路径修复完成"
