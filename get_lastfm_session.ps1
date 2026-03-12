$ConfigPath = ".\config.json"

if (!(Test-Path $ConfigPath)) {
    Write-Host "config.json not found"
    exit 1
}

$config = Get-Content $ConfigPath -Raw | ConvertFrom-Json

$apiKey = $config.lastfm_api_key
$secret = $config.lastfm_secret

if ([string]::IsNullOrEmpty($apiKey) -or [string]::IsNullOrEmpty($secret)) {
    Write-Host "Missing lastfm_api_key or lastfm_secret in config.json"
    exit 1
}

$apiUrl = "https://ws.audioscrobbler.com/2.0/"

Write-Host "Requesting auth token..."

$tokenResp = Invoke-RestMethod "$apiUrl?method=auth.getToken&api_key=$apiKey&format=json"
$token = $tokenResp.token

if (!$token) {
    Write-Host "Failed to obtain token"
    exit 1
}

$authUrl = "https://www.last.fm/api/auth/?api_key=$apiKey&token=$token"

Write-Host "Opening browser for authorization..."
Start-Process $authUrl

Write-Host ""
Write-Host "Authorize the app, then press ENTER..."
Read-Host

$signatureString = "api_key${apiKey}methodauth.getSessiontoken${token}${secret}"

$md5 = [System.Security.Cryptography.MD5]::Create()
$bytes = [System.Text.Encoding]::UTF8.GetBytes($signatureString)
$hash = $md5.ComputeHash($bytes)
$apiSig = ($hash | ForEach-Object { $_.ToString("x2") }) -join ""

Write-Host "Requesting session key..."

$body = @{
    method  = "auth.getSession"
    api_key = $apiKey
    token   = $token
    api_sig = $apiSig
    format  = "json"
}

$response = Invoke-RestMethod -Method Post -Uri $apiUrl -Body $body

$sessionKey = $response.session.key

if (!$sessionKey) {
    Write-Host "Failed to obtain session key"
    $response | ConvertTo-Json
    exit 1
}

$config.lastfm_session_key = $sessionKey

$config | ConvertTo-Json -Depth 10 | Set-Content $ConfigPath

Write-Host ""
Write-Host "Session key saved to config.json"
Write-Host "Session key: $sessionKey"