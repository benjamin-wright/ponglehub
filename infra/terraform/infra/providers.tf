provider "kubernetes" {
    config_context = "k3d-${var.cluster}"
}

provider "helm" {
    kubernetes {
        config_context = "k3d-${var.cluster}"
    }
}