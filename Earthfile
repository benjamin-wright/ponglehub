FROM scratch

all:
  BUILD ./services/auth-operator/operator+image
  BUILD ./services/auth-server+all
  BUILD ./services/db-init+image

repos-start:
  LOCALLY
  RUN ./infra/repos.sh start

repos-stop:
  LOCALLY
  RUN ./infra/repos.sh stop

libs:
  BUILD ./libraries/node/eslint-config-ponglehub+publish
  BUILD ./libraries/node/async+publish

generate:
  LOCALLY
  RUN rm -rf infra/manifests/generated
  RUN mkdir -p infra/manifests/generated
  RUN curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.10.1 sh -
  RUN helm template istio-operator ./istio-1.10.1/manifests/charts/istio-operator --namespace istio-operator > infra/manifests/generated/istio-operator.yaml
  RUN cp ./istio-1.10.1/manifests/charts/istio-operator/crds/* infra/manifests/generated/
  RUN rm -rf istio-1.10.1
  RUN curl -L https://github.com/knative/serving/releases/download/v0.22.0/serving-crds.yaml -o infra/manifests/generated/knative-serving-crds.yaml
  RUN curl -L https://github.com/knative/serving/releases/download/v0.22.0/serving-core.yaml -o infra/manifests/generated/knative-serving-core.yaml
  RUN curl -L https://github.com/knative/net-istio/releases/download/v0.22.0/net-istio.yaml -o infra/manifests/generated/knative-net-istio.yaml
  RUN yq eval-all \
        --inplace \
        '. |= (select(.kind=="Service" and .metadata.name=="knative-local-gateway") | .metadata.labels["experimental.istio.io/disable-gateway-port-translation"]="true")' \
        infra/manifests/generated/knative-net-istio.yaml


start:
  ARG REGISTRY_NAME=pongle-registry.localhost
  ARG REGISTRY_PORT=5000
  
  LOCALLY
  RUN k3d registry create $REGISTRY_NAME --port $REGISTRY_PORT
  RUN k3d cluster create pongle \
        --registry-use $REGISTRY_NAME \
        --k3s-server-arg "--no-deploy=traefik" \
        --kubeconfig-update-default=false \
        --volume $(pwd)/infra/manifests:/var/lib/rancher/k3s/server/manifests/preload \
        -p "80:80@loadbalancer" \
        --wait
  RUN mkdir -p .scratch
  RUN k3d kubeconfig get pongle > .scratch/kubeconfig

stop:
  ARG REGISTRY_NAME=pongle-registry.localhost

  LOCALLY
  RUN k3d cluster delete $CLUSTER_NAME || true
  RUN k3d registry delete $REGISTRY_NAME || true

clean:
  BUILD +stop
  BUILD +repos-stop

