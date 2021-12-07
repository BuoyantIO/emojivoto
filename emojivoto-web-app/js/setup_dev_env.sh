#!/bin/bash

# global vars
CONTAINER_ID=''

check_init_config() {
    if [[ -z "${AMBASSADOR_API_KEY}" ]]; then
        # you will need to set the AMBASSADOR_API_KEY via the command line
        # export AMBASSADOR_API_KEY='NTIyOWExZDktYTc5...'
        echo 'AMBASSADOR_API_KEY is not currently defined. Please set the environment variable in the shell e.g.'
        echo 'export AMBASSADOR_API_KEY=NTIyOWExZDktYTc5...'
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
    CONTAINER_ID=$(docker run -d -p8083:8083 -p8080:8080 --name ambassador-demo --cap-add=NET_ADMIN --device /dev/net/tun:/dev/net/tun --pull always --rm -it -e AMBASSADOR_API_KEY=$AMBASSADOR_API_KEY  -v ~/Library/Application\ Support:/root/.host_config  -v $(pwd):/opt/emojivoto/emojivoto-web-app/js datawire/emojivoto-node-and-go-demo )
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
    if [ ! command -v telepresence &> /dev/null ];  then
        echo "Installing Telepresence"
        sudo curl -fL https://app.getambassador.io/download/tel2/darwin/amd64/latest/telepresence -o /usr/local/bin/telepresence
        sudo chmod a+x /usr/local/bin/telepresence
    else
        echo "Telepresence already installed"
    fi    
}

connect_local_dev_env_to_remote() {
    export KUBECONFIG=./emojivoto_k8s_context.yaml
    echo 'Connecting local dev env to remote K8s cluster'
    telepresence intercept web-app-57bc7c4959 --service web-app --port 8083:80
}

open_editor() {
    echo 'Opening editor'

    # replace this line with your editor of choice, e.g. VS code, Intelli J
    code .
}

display_instructions_to_user () {
    echo ''
    echo 'INSTRUCTIONS FOR DEVELOPMENT'
    echo '============================'
    echo 'To set the correct Kubernetes context on this shell, please execute:'
    echo 'export KUBECONFIG=./emojivoto_k8s_context.yaml'
}

check_init_config
run_dev_container
connect_to_k8s
install_telepresence
connect_local_dev_env_to_remote
open_editor
display_instructions_to_user

# happy coding!
