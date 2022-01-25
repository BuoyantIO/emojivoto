#!/bin/bash

# global vars
CONTAINER_ID=''
OS=''
ARCH=''
MOUNT_VOLUME=''

has_cli() {

    hasCurl=$(which curl)
    if [ "$?" = "1" ]; then
        echo "You need curl to use this script."
        exit 1
    fi
    hasKubectl=$(which kubectl)
    if [ "$?" = "1" ]; then
        echo "You need kubectl to use this script. https://kubernetes.io/docs/tasks/tools/#kubectl"
        exit 1
    fi
}

set_os_arch() {
    uname=$(uname)

    case $uname in
        "Darwin")
            OS="darwin"
            MOUNT_VOLUME_LOCAL=~/Library/Application Support
            ;;
        "Linux")
            OS="linux"
            MOUNT_VOLUME_LOCAL=$(if [[ "$XDG_CONFIG_HOME" ]]; then echo "$XDG_CONFIG_HOME"; else echo "$HOME/.config"; fi)
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
    echo 'Running dev container (and downloading if necessary)'

    # check if dev container is already running and kill if so
    CONTAINER_ID=$(docker inspect --format="{{.Id}}" "ambassador-demo" )
    if [ ! -z "$CONTAINER_ID" ]; then
        docker kill $CONTAINER_ID
    fi

    # run the dev container, exposing 8081 gRPC port and volume mounting code directory
    CONTAINER_ID=$(docker run -d -p8083:8083 -p8080:8080 --name ambassador-demo --cap-add=NET_ADMIN --device /dev/net/tun:/dev/net/tun --pull always --rm -it -e AMBASSADOR_API_KEY=$AMBASSADOR_API_KEY  -v $MOUNT_VOLUME_LOCAL:/root/.host_config  -v $(pwd):/opt/emojivoto/emojivoto-web-app/js datawire/emojivoto-node-and-go-demo )
}

connect_to_k8s() {
    echo 'Extracting KUBECONFIG from container and connecting to cluster'
    until docker cp $CONTAINER_ID:/opt/telepresence-demo-cluster.yaml ./emojivoto_k8s_context.yaml > /dev/null 2>&1; do
        echo '.'
        sleep 1s
    done

    export KUBECONFIG=./emojivoto_k8s_context.yaml

    echo 'Connected to cluster. Listing services in default namespace'
    kubectl get svc
}

install_telepresence() {
    echo 'Configuring Telepresence'
    hasTelepresence=$(which telepresence)
    if [ "$?" = "1" ]; then
        echo "Installing Telepresence"
        sudo curl -fL https://app.getambassador.io/download/tel2/${OS}/${ARCH}/latest/telepresence -o /usr/local/bin/telepresence
        sudo chmod a+x /usr/local/bin/telepresence
    else
        echo "Telepresence already installed"
    fi    
}

connect_local_dev_env_to_remote() {
    export KUBECONFIG=./emojivoto_k8s_context.yaml
    echo 'Connecting local dev env to remote K8s cluster'
    telepresence login --apikey=$AMBASSADOR_API_KEY
    telepresence intercept web-app-57bc7c4959 -n emojivoto --service web-app --port 8083:80 --ingress-port 80 --ingress-host ambassador.ambassador --ingress-l5 ambassador.ambassador
    telOut=$?
    if [ $telOut != 0 ]; then
        exit $telOut
    fi
}

open_editor() {
    echo 'Opening editor'
    # let the user see the output before opening the editor
    sleep 2

    # replace this line with your editor of choice, e.g. VS code, Intelli J
    $EDITOR .
}

display_instructions_to_user () {
    echo ''
    echo 'INSTRUCTIONS FOR DEVELOPMENT'
    echo '============================'
    echo 'To set the correct Kubernetes context on this shell, please execute:'
    echo 'export KUBECONFIG=./emojivoto_k8s_context.yaml'
}

has_cli
set_os_arch
check_init_config
run_dev_container
connect_to_k8s
install_telepresence
connect_local_dev_env_to_remote
open_editor
display_instructions_to_user

# happy coding!
