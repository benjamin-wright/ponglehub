locals {
    scratch_dir = "${abspath(path.module)}/.scratch"
}

resource "docker_network" "registry" {
    name = "${var.cluster}-network"
}

resource "docker_container" "registry" {
    image = "sha256:2d4f4b5309b1e41b4f83ae59b44df6d673ef44433c734b14c1c103ebca82c116"
    name = "${var.cluster}-registry"
    restart = "always"
    ports {
        internal = 5000
        external = var.registry_port
    }
    networks_advanced {
        name = docker_network.registry.name
    }
}

resource "null_resource" "cluster" {
    triggers = {
        SCRATCH_DIR = local.scratch_dir
        CLUSTER_NAME = var.cluster
        NETWORK_NAME = docker_network.registry.name
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
