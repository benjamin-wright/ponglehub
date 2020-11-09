resource "kubernetes_namespace" "infra" {
  metadata {
    name = "infra"
    annotations = {
      "linkerd.io/inject" = "enabled"
    }
  }

  # depends_on = [
  #   helm_release.linkerd
  # ]
}

resource "tls_private_key" "ingress_ca" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_self_signed_cert" "ingress_ca" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.ingress_ca.private_key_pem

  subject {
    common_name  = var.domain
    organization = "Ponglehub, Inc"
  }

  validity_period_hours = 175200

  allowed_uses = [
    "cert_signing",
    "ca_signing"
  ]

  is_ca_certificate = true
}

resource "tls_private_key" "ingress" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_cert_request" "ingress" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.ingress.private_key_pem

  subject {
    common_name  = var.domain
    organization = "Ponglehub, Inc"
  }

  dns_names = [
    "*.${var.domain}"
  ]
}

resource "tls_locally_signed_cert" "ingress" {
  cert_request_pem   = tls_cert_request.ingress.cert_request_pem
  ca_key_algorithm   = "ECDSA"
  ca_private_key_pem = tls_private_key.ingress_ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ingress_ca.cert_pem

  validity_period_hours = 3600

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth"
  ]
}

resource "local_file" "ingress_key" {
  content = tls_private_key.ingress.private_key_pem
  filename = "${local.scratch_dir}/ingress.key"
}

resource "local_file" "ingress_cert" {
  content = tls_locally_signed_cert.ingress.cert_pem
  filename = "${local.scratch_dir}/ingress.crt"
}

resource "local_file" "ca_cert" {
  content = tls_self_signed_cert.ingress_ca.cert_pem
  filename = "${local.scratch_dir}/ingress-ca.crt"
}

resource "kubernetes_secret" "infra_tls_secret" {
  metadata {
    name = "tls-secret"
    namespace = kubernetes_namespace.infra.metadata[0].name
  }

  data = {
    "tls.key" = tls_private_key.ingress.private_key_pem
    "tls.crt" = tls_locally_signed_cert.ingress.cert_pem
  }

  type = "Opaque"
}