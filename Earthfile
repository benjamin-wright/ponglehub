FROM scratch

all:
  BUILD ./services/auth-operator/operator+image
  BUILD ./services/auth-server+all
  BUILD ./services/db-init+image

repos:
  LOCALLY
  RUN ./infra/repos.sh
  RUN helm dep update helm/tests

libs:
  BUILD ./libraries/node/eslint-config-ponglehub+publish
  BUILD ./libraries/node/async+publish

repos-stop:
  LOCALLY
  RUN ./infra/repos-stop.sh

generate:
  LOCALLY
  # RUN curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.10.1 sh -
  # RUN helm template istio-operator ./istio-1.10.1/manifests/charts/istio-operator --namespace istio-operator > infra/manifests/istio-operator.yaml
  # RUN rm -rf istio-1.10.1
  RUN curl -L https://github.com/knative/serving/releases/download/v0.22.0/serving-crds.yaml -o infra/manifests/knative-serving-crds.yaml
  RUN curl -L https://github.com/knative/serving/releases/download/v0.22.0/serving-core.yaml -o infra/manifests/knative-serving-core.yaml
  RUN curl -L https://github.com/knative/net-istio/releases/download/v0.22.0/net-istio.yaml -o infra/manifests/knative-net-istio.yaml

infra:
  LOCALLY
  RUN ./infra/start.sh

infra-stop:
  LOCALLY
  RUN ./infra/stop.sh

start:
  BUILD +repos
  BUILD +infra

stop:
  BUILD +infra-stop

clean:
  BUILD +infra-stop
  BUILD +repos-stop

