resource "helm_release" "kafka" {
  name = "kafka"
  namespace = kubernetes_namespace.infra.metadata[0].name
  chart = "strimzi-kafka-operator"
  repository = "https://strimzi.io/charts/"
  version = "0.20.0"

  values = [
    file("values/kafka.yaml")
  ]

  depends_on = [
    kubernetes_namespace.ponglehub
  ]
}

resource "helm_release" "cockroach" {
  name = "cockroach"
  namespace = kubernetes_namespace.infra.metadata[0].name
  chart = "cockroachdb"
  repository = "https://charts.cockroachdb.com/"
  version = "4.1.10"

  values = [
    file("values/cockroach.yaml")
  ]
}
