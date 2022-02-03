#!/bin/bash

# global vars
CONTAINER_ID=''
OS=''
ARCH=''
MOUNT_VOLUME_LOCAL=''
USE_TELEMETRY=
OPEN_EDITOR=
ACTIVITY_REPORT_TYPE='INTERMEDIATE_CLOUD_TOUR_SCRIPT'
EMOJIVOTO_NS='emojivoto'
TOTAL_STEPS='7'


use_telemetry() {
    USE_TELEMETRY=true
}

send_telemetry() {
    if [ $USE_TELEMETRY = true ]; then
        action=$1
        ambassador_cloud_url="https://auth.datawire.io"
        application_activities_url="${ambassador_cloud_url}/api/applicationactivities"
        curl -X POST \
          -H "X-Ambassador-API-Key: $AMBASSADOR_API_KEY" \
          -H "Content-Type: application/json" \
          -d '{"type": "'$ACTIVITY_REPORT_TYPE'", "extraProperties": {"action":"'"$action"'","os":"'"$OS"'","arch":"'"$ARCH"'"}}' \
          -s \
          $application_activities_url > /dev/null 2>&1
    fi
}

display_step() {
    echo -n "*** Step $1/$TOTAL_STEPS: "
}

has_cli() {
    display_step 1
    meet_requirements=true
    echo 'Checking required tools ... '
    _=$(which curl)
    if [ "$?" = "1" ]; then
        echo "You need curl to use this script."
        meet_requirements=false
    fi
    _=$(which kubectl)
    if [ "$?" = "1" ]; then
        echo "You need kubectl to use this script. https://kubernetes.io/docs/tasks/tools/#kubectl"
        meet_requirements=false
    fi
    _=$(which docker)
    if [ "$?" = "1" ]; then
        echo "You need docker to use this script. https://docs.docker.com/engine/install/"
        meet_requirements=false
    fi

    if (! docker stats --no-stream &> /dev/null); then
        echo "Docker daemon is not running"
        meet_requirements=false
    fi
    if [ $meet_requirements = false ]; then
        exit 1
    fi
    echo 'Complete.'
}

set_os_arch() {
    uname=$(uname)

    case $uname in
        "Darwin")
            OS="darwin"
            MOUNT_VOLUME_LOCAL=~/Library/Application\ Support
            OPEN_EDITOR=open
            ;;
        "Linux")
            OS="linux"
            MOUNT_VOLUME_LOCAL=$(if [[ "$XDG_CONFIG_HOME" ]]; then echo "$XDG_CONFIG_HOME"; else echo "$HOME/.config"; fi)
            OPEN_EDITOR=xdg-open
            ;;
        *)
            fatal "Unsupported os $uname"
    esac

    if [ -z "$ARCH" ]; then
        ARCH=$(uname -m)
    fi
    case $ARCH in
        amd64)
            ARCH=amd64
            ;;
        x86_64)
            ARCH=amd64
            ;;
        arm64)
            ARCH=arm64
            ;;
        aarch64)
            ARCH=arm64
            ;;
        *)
            fatal "Unsupported architecture $ARCH"
    esac
}

check_init_config() {
    display_step 2
    echo 'Checking for AMBASSADOR_API_KEY environment variable'
    if [[ -z "${AMBASSADOR_API_KEY}" ]]; then
        # you will need to set the AMBASSADOR_API_KEY via the command line
        # export AMBASSADOR_API_KEY='NTIyOWExZDktYTc5...'
        echo 'AMBASSADOR_API_KEY is not currently defined. Please set the environment variable in the shell e.g.'
        echo 'export AMBASSADOR_API_KEY=NTIyOWExZDktYTc5...'
        echo 'You can get an AMBASSADOR_API_KEY and free remote demo cluster by taking the tour of Ambassador Cloud at https://app.getambassador.io/cloud/services?openCloudTour=true '
        echo 'During the tour be sure to copy the AMBASSADOR_API_KEY from the "docker run" command'
        exit
    fi
}

run_dev_container() {
    display_step 4
    echo 'Configuring development container. This container encapsulates all the dependencies needed to run the emojivoto-web-app locally.'
    echo 'This may take a few moments to download and start.'

    # check if dev container is already running and kill if so
    CONTAINER_ID=$(docker inspect --format="{{.Id}}" "ambassador-demo" 2>/dev/null)
    if [ -n "$CONTAINER_ID" ]; then
        _=$(docker kill "$CONTAINER_ID")
    fi

    # run the dev container, exposing 8081 gRPC port and volume mounting code directory
    docker run -d -p8083:8083 -p8080:8080 --name ambassador-demo --cap-add=NET_ADMIN --device /dev/net/tun:/dev/net/tun --pull always --rm -it -e AMBASSADOR_API_KEY=$AMBASSADOR_API_KEY  -v "${MOUNT_VOLUME_LOCAL}":/root/.host_config  -v $(pwd):/opt/emojivoto/emojivoto-web-app/js datawire/emojivoto-node-and-go-demo
    CONTAINER_ID=$(docker ps --filter 'name=ambassador-demo' --format '{{.ID}}')
    send_telemetry "devContainerStarted"    
}

connect_to_k8s() {
    display_step 5
    echo 'Extracting KUBECONFIG from container'
    until docker cp $CONTAINER_ID:/opt/telepresence-demo-cluster.yaml ./emojivoto_k8s_context.yaml > /dev/null 2>&1; do
        echo '.'
        sleep 2
    done

    export KUBECONFIG=./emojivoto_k8s_context.yaml

    echo "Listing services in ${EMOJIVOTO_NS} namespace"
    listSVC=$(kubectl --namespace ${EMOJIVOTO_NS} get svc)
    echo "$listSVC"
    send_telemetry "connectedToK8S"    
}

install_upgrade_telepresence() {
    display_step 3
    install_telepresence=false
    echo -n 'Checking for Telepresence ... '
    _=$(which telepresence)
    if [ "$?" = "1" ]; then
        install_telepresence=true
        echo "Installing Telepresence"
    else
        if telepresence version|grep upgrade >/dev/null 2>&1; then
            install_telepresence=true
            # daemon and client need to have same version
            _=$(telepresence quit)
            echo "Upgrading Telepresence"
        else
            echo "Telepresence already installed"
            send_telemetry "telepresenceAlreadyInstalled"
        fi
    fi    
    if [ $install_telepresence = true ]; then
        sudo curl -fL https://app.getambassador.io/download/tel2/${OS}/${ARCH}/latest/telepresence -o /usr/local/bin/telepresence
        sudo chmod a+x /usr/local/bin/telepresence
        send_telemetry "telepresenceInstalled"
    fi
}

connect_local_dev_env_to_remote() {
    export KUBECONFIG=./emojivoto_k8s_context.yaml
    display_step 6
    echo 'Connecting local development environment to remote K8s cluster'
    telepresence login --apikey=${AMBASSADOR_API_KEY}
    telepresence intercept web-app-57bc7c4959 -n ${EMOJIVOTO_NS} --service web-app --port 8083:80 --ingress-port 80 --ingress-host ambassador.ambassador --ingress-l5 ambassador.ambassador
    telOut=$?
    if [ $telOut != 0 ]; then
        send_telemetry "interceptFailed"
        exit $telOut
    fi
    send_telemetry "interceptCreated"
}

open_editor() {
    display_step 7
    echo 'Opening editor'
    # let the user see the output before opening the editor
    sleep 2

    # replace this line with your editor of choice, e.g. VS code, Intelli J
    $OPEN_EDITOR components/Vote.jsx
}

display_instructions_to_user () {
    echo ''
    echo 'INSTRUCTIONS FOR DEVELOPMENT'
    echo '============================'
    echo 'To set the correct Kubernetes context on this shell, please execute:'
    echo 'export KUBECONFIG=./emojivoto_k8s_context.yaml'
}

use_telemetry
has_cli
set_os_arch
check_init_config
install_upgrade_telepresence
run_dev_container
connect_to_k8s
connect_local_dev_env_to_remote
open_editor
display_instructions_to_user

# happy coding!
