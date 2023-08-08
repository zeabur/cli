# get the latest release of zeabur cli
$release_url = "https://api.github.com/repos/zeabur/cli/releases"
$tag = (Invoke-WebRequest -Uri $release_url -UseBasicParsing | ConvertFrom-Json)[0].tag_name
$loc = "$HOME\AppData\Local\zeabur-cli"
$url = ""
$arch = $env:PROCESSOR_ARCHITECTURE
$release = $tag.Trim("v")
$releases_api_url = "https://github.com/zeabur/cli/releases/download/$tag/zeabur_${release}_windows"

if ($arch -eq "x86") {
    $url = "${releases_api_url}_386.exe"
} elseif ($arch -eq "arm64") {
    $url = "${releases_api_url}_arm64.exe"
} else {
    $url = "${releases_api_url}_amd64.exe"
}

if (Test-Path -path $loc) {
    Remove-Item $loc -Recurse -Force
}

Write-Host "Installing zeabur cli version $tag" -ForegroundColor DarkCyan

Invoke-WebRequest $url -outfile zeabur.exe

New-Item -ItemType "directory" -Path $loc

Move-Item -Path zeabur.exe -Destination $loc

[System.Environment]::SetEnvironmentVariable("Path", $Env:Path + ";$loc", [System.EnvironmentVariableTarget]::User)

if (Test-Path -path $loc) {
    Write-Host "Thanks for installing Zeabur CLI! Now Refresh your powershell" -ForegroundColor DarkGreen
    Write-Host "If this is your first time using the CLI, be sure to run 'zeabur --help' first." -ForegroundColor DarkGreen
} else {
    Write-Host "Download failed" -ForegroundColor Red
    Write-Host "Please try again later" -ForegroundColor Red
}
