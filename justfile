cluster_name := "pongle"
registry_name := "pongle-registry.localhost"
registry_port := "5000"

start: create-cluster setup-context wait-for-traefik
stop: delete-cluster clear-context

create-cluster:
    #!/usr/bin/env bash
    set -euxo pipefail

    if ! k3d cluster list | grep -qw {{ cluster_name }}; then
        k3d cluster create {{ cluster_name }} \
            --registry-create {{ registry_name }}:0.0.0.0:{{ registry_port }} \
            --kubeconfig-update-default=false \
            -p "80:80@loadbalancer" \
            --wait;
    else
        echo "cluster {{ cluster_name }} already exists!"
    fi

setup-context:
    @mkdir -p .scratch
    @k3d kubeconfig get {{ cluster_name }} > .scratch/kubeconfig
    @KUBECONFIG=.scratch/kubeconfig kubectl config use-context k3d-{{ cluster_name }}

delete-cluster:
    if k3d cluster list | grep -qw {{ cluster_name }}; then \
        k3d cluster delete {{ cluster_name }}; \
    fi

clear-context:
    if [[ -f .scratch/kubeconfig ]]; then \
        rm .scratch/kubeconfig; \
    fi

wait-for-traefik:
    #!/usr/bin/env bash
    LAST_STATUS=""
    STATUS=""
    
    echo "Waiting for traefik to start..."

    while [[ "$STATUS" != "Running" ]]; do
        sleep 1
        STATUS=$(kubectl get pods -n kube-system -o json | jq '.items[] | select(.metadata.name | startswith("traefik")) | .status.phase' -r)
        if [[ "$STATUS" != "$LAST_STATUS" ]]; then
            echo "traefik pod is '$STATUS'"
        fi
        LAST_STATUS="$STATUS"
    done

    echo "done"