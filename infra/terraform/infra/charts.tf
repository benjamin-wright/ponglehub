resource "helm_release" "verdaccio" {
  name = "verdaccio"
  namespace = kubernetes_namespace.infra.metadata[0].name
  chart = "verdaccio"
  repository = "https://charts.verdaccio.org"
  version = "0.10.0"

  values = [
    "${file("values/verdaccio.yaml")}"
  ]
}

resource "helm_release" "chartmuseum" {
  name = "chartmuseum"
  namespace = kubernetes_namespace.infra.metadata[0].name
  chart = "chartmuseum"
  repository = "https://kubernetes-charts.storage.googleapis.com/"
  version = "2.14.0"

  values = [
    "${file("values/chartmuseum.yaml")}"
  ]
}

resource "helm_release" "kafka" {
  name = "kafka"
  namespace = kubernetes_namespace.infra.metadata[0].name
  chart = "strimzi-kafka-operator"
  repository = "https://strimzi.io/charts/"
  version = "0.19.0"
}

resource "helm_release" "cockroach" {
  name = "cockroach"
  namespace = kubernetes_namespace.infra.metadata[0].name
  chart = "cockroachdb"
  repository = "https://charts.cockroachdb.com/"
  version = "4.1.10"

  values = [
    "${file("values/cockroach.yaml")}"
  ]
}
