
resource "tls_private_key" "linkerd_ca" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_self_signed_cert" "linkerd_ca" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.linkerd_ca.private_key_pem

  subject {
    common_name  = "identity.linkerd.cluster.local"
    organization = "Ponglehub, Inc"
  }

  validity_period_hours = 175200

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "cert_signing",
    "ca_signing"
  ]

  is_ca_certificate = true
}

resource "tls_private_key" "linkerd" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_cert_request" "linkerd" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.linkerd.private_key_pem

  subject {
    common_name  = "identity.linkerd.cluster.local"
    organization = "Ponglehub, Inc"
  }
}

resource "tls_locally_signed_cert" "linkerd" {
  cert_request_pem   = tls_cert_request.linkerd.cert_request_pem
  ca_key_algorithm   = "ECDSA"
  ca_private_key_pem = tls_private_key.linkerd_ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.linkerd_ca.cert_pem

  validity_period_hours = 3600

  is_ca_certificate = true

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "cert_signing",
    "ca_signing"
  ]
}

resource "helm_release" "linkerd" {
  name = "linkerd"
  chart = "linkerd2"
  repository = "https://helm.linkerd.io/edge"

  set {
    name = "grafana.enabled"
    value = false
  }

  set {
    name = "ingress.hostname"
    value = "linkerd.ponglehub.co.uk"
  }

  set {
    name = "enforcedHostRegexp"
    value = "^linkerd\\.ponglehub\\.co\\.uk$"
  }

  set {
    name = "identity.issuer.crtExpiry"
    value = timeadd(timestamp(), "175200h")
  }

  set {
    name = "global.identityTrustAnchorsPEM"
    value = tls_self_signed_cert.linkerd_ca.cert_pem
  }

  set {
    name = "identity.issuer.tls.crtPEM"
    value = tls_locally_signed_cert.linkerd.cert_pem
  }

  set {
    name = "identity.issuer.tls.keyPEM"
    value = tls_private_key.linkerd.private_key_pem
  }

  lifecycle {
    ignore_changes = [
      set,
    ]
  }
}

resource "kubernetes_secret" "linkerd_ssl_secret" {
  depends_on = [
    helm_release.linkerd
  ]

  metadata {
    name = "tls-secret"
    namespace = "linkerd"
  }

  data = {
    "tls.key" = tls_private_key.ingress.private_key_pem
    "tls.crt" = tls_locally_signed_cert.ingress.cert_pem
  }

  type = "Opaque"
}

resource "kubernetes_ingress" "linkerd_ingress" {
  depends_on = [
    helm_release.linkerd
  ]

  metadata {
    name = "web-ingress"
    namespace = "linkerd"
    annotations = {
      "ingress.kubernetes.io/custom-request-headers" = "l5d-dst-override:linkerd-web.linkerd.svc.cluster.local:8084"
    }
  }

  spec {
    rule {
      host = "linkerd.${var.domain}"
      http {
        path {
          backend {
            service_name = "linkerd-web"
            service_port = 8084
          }
        }
      }
    }

    tls {
      hosts = [ "linkerd.${var.domain}" ]
      secret_name = "tls-secret"
    }
  }
}