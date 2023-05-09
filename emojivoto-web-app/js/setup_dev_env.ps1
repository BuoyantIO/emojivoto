param([string]$action="TP_SCIPT")

$Global:USE_TELEMETRY=$true
$Global:ACTION=$action
$Global:ACTIVITY_REPORT_TYPE='INTERMEDIATE_CLOUD_TOUR_SCRIPT'
$Global:OS_INFO=Get-ComputerInfo
$Global:SEGMENT=$action
$Global:TOTAL_STEPS=7

function use_telemetry {
    $Global:USE_TELEMETRY=$true
}

function send_telemetry{
    if ($Global:USE_TELEMETRY) {
        [string]$ambassador_cloud_url="https://auth.datawire.io"
        [string]$application_activities_url="$ambassador_cloud_url/api/applicationactivities"
        [string]$body = @{
            type = "$Global:ACTIVITY_REPORT_TYPE"
            extraProperties = @{
                "action" = "$Global:ACTION"
                "os" = "$Global:OS_INFO.OperatingSystem"
                "arch" = "$Global:OS_INFO.OperatingSystemArchitecture"
                "segment" = "$Global:SEGMENT"
            }
        } | ConvertTo-Json
        $headers=@{
            "X-Ambassador-API-KEY" = "$Env:AMBASSADOR_API_KEY"
            "Content-Type" = "application/json"
        }
        Invoke-WebRequest -Uri $application_activities_url -Headers $headers -Body $body -Method 'POST'
    }
}

function display_step{
    param(
        [Int32]$step_number
    )
    Write-Host "$step_number/$Global:TOTAL_STEPS"
}

function has_cli() {
    display_step(1)
    $meet_requirements=$true
    Write-Host "Checking required tools ..."
    if(-not (Get-Command curl -ErrorAction SilentlyContinue)){
        Write-Host "You need curl to use this script."
        $meet_requirements=$false
    }
    if(-not (Get-Command kubectl -ErrorAction SilentlyContinue)){
        Write-Host "You need kubectl to use this script. https://kubernetes.io/docs/tasks/tools/#kubectl"
        $meet_requirements=$false
    }
    if(-not (Get-Command docker -ErrorAction SilentlyContinue)){
        Write-Host "You need docker to use this script. https://docs.docker.com/engine/install/"
        $meet_requirements=$false
    }
    if (-not (docker stats --no-stream 2> $null)) {
        Write-Host "El demonio de Docker no está ejecutándose"
        $meet_requirements = $false
    }
    
}

use_telemetry
send_telemetry
has_cli
