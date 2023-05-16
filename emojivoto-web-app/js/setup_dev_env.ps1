param([string]$action="TP_SCIPT")

$Global:USE_TELEMETRY=$true
$Global:ACTION=$action
$Global:ACTIVITY_REPORT_TYPE='INTERMEDIATE_CLOUD_TOUR_SCRIPT'
$Global:OS_INFO=Get-ComputerInfo
$Global:SEGMENT=$action
$Global:TOTAL_STEPS=7
$Global:OS="windows"
$Global:ARCH="AMD64"
$Global:EMOJIVOTO_NS="emojivoto"
$Env:AMBASSADOR_API_KEY="ZjRlMTAzMDYtYjU5Ni00NTU3LTk2YjgtMTM3MTMzOWZjNmQxOlB0T3A0VFg0TG5JZVVucFJ1VmwzNEJOMHRFNk1nakRCTVIwSw=="

function use_telemetry {
    $Global:USE_TELEMETRY=$true
}

function send_telemetry{
    param(
        [string]$action
    )
    if ($Global:USE_TELEMETRY) {
        [string]$ambassador_cloud_url="https://auth.datawire.io"
        [string]$application_activities_url="$ambassador_cloud_url/api/applicationactivities"
        [string]$body = @{
            type = "$Global:ACTIVITY_REPORT_TYPE"
            extraProperties = @{
                "action" = "$action"
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

function has_cli {
    display_step(1)
    $meet_requirements=$true
    Write-Host "Checking required tools ..."
    if(-not (Get-Command Invoke-WebRequest -ErrorAction SilentlyContinue)){
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
        $meet_requirements=$false
    }
    if(-not $meet_requirements){
        exit
    }
    Write-Host "Complete."
}

function set_os_arch{
    if ($Env:PROCESSOR_ARCHITECTURE -eq "AMD64") {
        $Global:ARCH=$Env:PROCESSOR_ARCHITECTURE
    } else {
        Write-Host "Unsupported architecture $Env:PROCESSOR_ARCHITECTURE"
        exit 
    }
    
}

function check_init_config{
    display_step(2)
    Write-Host "Checking for AMBASSADOR_API_KEY environment variable"
    if(-not ($Env:AMBASSADOR_API_KEY)){
        # you will need to set the AMBASSADOR_API_KEY via the command line
        # New-Item -Path Env:\AMBASSADOR_API_KEY -Value 'NTIyOWExZDktYTc5...'
        Write-Host 'AMBASSADOR_API_KEY is not currently defined. Please set the environment variable in the shell e.g.'
        Write-Host 'New-Item -Path Env:\AMBASSADOR_API_KEY -Value NTIyOWExZDktYTc5...'
        Write-Host 'You can get an AMBASSADOR_API_KEY and free remote demo cluster by taking the tour of Ambassador Cloud at https://app.getambassador.io/cloud/welcome?tour=intermediate '
        Write-Host 'During the tour be sure to copy the AMBASSADOR_API_KEY from the "docker run" command'
        exit
    }
}

function run_dev_container{
    display_step(6)
    Write-Host 'Configuring development container. This container encapsulates all the dependencies needed to run the emojivoto-web-app locally.'
    Write-Host 'This may take a few moments to download and start.'

    # check if dev container is already running and kill if so
    $CONTAINER_ID=$(docker inspect --format='{{.Id}}' 'ambassador-demo' 2>$null)
    if ($CONTAINER_ID) {
        docker kill $CONTAINER_ID >$null
    }
    docker run -d --name ambassador-demo --pull always --network=container:tp-default --rm -it -v ${Get-Location}:\opt\emojivoto\emojivoto-web-app\js datawire/intermediate-tour
    $CONTAINER_ID=$(docker ps --filter 'name=ambassador-demo' --format '{{.ID}}')
    send_telemetry("devContainerStarted")
}

function connect_to_k8s{
    display_step(4)
    Write-Host 'Getting KUBECONFIG from demo cluster'

    $demo_cluster_url = "https://auth.datawire.io/api/democlusters/telepresence-demo/config"
    if ($env:SYSTEMA_ENV -eq "staging") {
        $demo_cluster_url = "https://staging-auth.datawire.io/api/democlusters/telepresence-demo/config"
    }

    Invoke-WebRequest $demo_cluster_url -Headers @{ "X-Ambassador-API-Key" = $env:AMBASSADOR_API_KEY } -ContentType "application/json" -OutFile 'emojivoto_k8s_context.yaml'
    $env:KUBECONFIG = './emojivoto_k8s_context.yaml'
    kubectl config set-context --current --namespace=emojivoto

    Write-Host "Listing services in ${Global:EMOJIVOTO_NS} namespace"
    $listSVC = kubectl --namespace ${Global:EMOJIVOTO_NS} get svc
    $listSVC
    send_telemetry "connectedToK8S"
}

function install_upgrade_telepresence{
    display_step(3)
    $install_telepresence = $false
    Write-Host -NoNewline 'Checking for Telepresence ... '
    if (!(Get-Command telepresence -ErrorAction SilentlyContinue)) {
        $install_telepresence = $true
        Write-Host "Installing Telepresence"
    } else {
        # Ensure that running telepresence daemons are stopped. A running daemon
        # has its current working directory set to the directory where it was first
        # started, and since the KUBECONFIG is set to a relative directory in this
        # script, a previously started daemon might resolve it incorrectly.
        telepresence quit -s > $null
        if ((telepresence version) -match "upgrade") {
            $install_telepresence = $true
            Write-Host "Upgrading Telepresence"
        } else {
            Write-Host "Telepresence already installed"
            send_telemetry("telepresenceAlreadyInstalled")
        }
    }
    #if ($install_telepresence) {
        $telepresence_download_url = "https://app.getambassador.io/download/tel2/windows/amd64/latest/telepresence-setup.exe"
        #Invoke-WebRequest $telepresence_download_url -OutFile telepresence.zip
        #Invoke-WebRequest $telepresence_download_url -OutFile telepresence-setup.exe
        $telepresenceOutput = Start-Process .\telepresence-setup.exe -NoNewWindow -Wait
    if ($telepresenceOutput.ExitCode -ne 0) {
        Write-Host "You need to install telepresence to continue with the execution of this script."
        exit 1
    }
        #powershell.exe -Confirm:$false -ExecutionPolicy bypass -c " . '.\telepresence-setup.exe';"
        #Expand-Archive -Path telepresence.zip -DestinationPath telepresenceInstaller/telepresence
        #Remove-Item 'telepresence.zip'
        #Set-Location telepresenceInstaller/telepresence
        #powershell.exe -ExecutionPolicy bypass -c " . '.\install-telepresence.ps1';"        
        #Set-Location ../..
        #Remove-Item telepresenceInstaller -Recurse -Confirm:$false -Force
        #send_telemetry("telepresenceInstalled")
    #}
}

function connect_local_dev_env_to_remote{
    $env:KUBECONFIG = './emojivoto_k8s_context.yaml'
    display_step 5
    Write-Host 'Connecting local development environment to remote K8s cluster'

    $svcName = "ambassador"
    kubectl get svc ambassador -n ambassador
    $ambassadorSvcOut = $LASTEXITCODE
    if ($ambassadorSvcOut -ne 0) {
        $svcName = "edge-stack"
    }

    telepresence quit > $null
    telepresence helm upgrade --team-mode > $null
    telepresence login --apikey="$Env:AMBASSADOR_API_KEY"
    telepresence quit -s > $null
    telepresence connect --docker > $null

    $interceptName = (kubectl get rs -n emojivoto --selector=app=web-app --no-headers -o custom-columns=":metadata.name")
    telepresence intercept "$interceptName" --docker --context default -n "$Global:EMOJIVOTO_NS" --service web-app --port 8083:80 --ingress-port 80 --ingress-host "$svcName.ambassador" --ingress-l5 "$svcName.ambassador"

    $telOut = $LASTEXITCODE
    if ($telOut -ne 0) {
        send_telemetry "interceptFailed"
        exit $telOut
    }
    send_telemetry "interceptCreated"
}

function open_editor() {
    display_step(7)
    Write-Host 'Opening editor'
    # let the user see the output before opening the editor
    Start-Sleep 2

    # replace this line with your editor of choice, e.g. VS code, Intelli J
    Invoke-Item components\Vote.jsx
}

function display_instructions_to_user () {
    Write-Host ''
    Write-Host 'INSTRUCTIONS FOR DEVELOPMENT'
    Write-Host '============================'
    Write-Host 'To set the correct Kubernetes context on this shell, please execute:'
    Write-Host 'export KUBECONFIG=./emojivoto_k8s_context.yaml'
}

use_telemetry
has_cli
set_os_arch
check_init_config
install_upgrade_telepresence
connect_to_k8s
connect_local_dev_env_to_remote
run_dev_container
open_editor
display_instructions_to_user

# happy coding!
