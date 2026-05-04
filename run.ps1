param (
    [Parameter(Mandatory=$false)]
    [string]$XXXX = "8080"
)

$env:XXXX = $XXXX
Write-Host "Starting application on port $XXXX..." -ForegroundColor Green
docker-compose up --build