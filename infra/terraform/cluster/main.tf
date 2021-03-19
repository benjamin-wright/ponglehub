locals {
    scratch_dir = "${abspath(path.module)}/.scratch"
}

resource "null_resource" "cluster" {
    triggers = {
        SCRATCH_DIR = local.scratch_dir
        CLUSTER_NAME = var.cluster
        NETWORK_NAME = "${var.cluster}-network"
        REGISTRY_PORT = var.registry_port
        WORKER_NODES = 0
    }

    provisioner "local-exec" {
        command = "./scripts/start-cluster.sh"
        environment = self.triggers
    }

    provisioner "local-exec" {
        when = destroy
        command = "./scripts/stop-cluster.sh"
        environment = self.triggers
    }
}
