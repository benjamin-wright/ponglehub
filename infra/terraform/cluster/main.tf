locals {
    scratch_dir = "${abspath(path.module)}/.scratch"
}

resource "docker_network" "registry" {
    name = "${var.cluster}-network"
}

data "docker_registry_image" "registry" {
  name = "registry:2"
}

resource "docker_image" "registry" {
  name          = data.docker_registry_image.registry.name
  pull_triggers = [data.docker_registry_image.registry.sha256_digest]
}

resource "docker_container" "registry" {
    image = docker_image.registry.latest
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

resource "null_resource" "geppetto" {
    provisioner "local-exec" {
        command = "cd ../../../tools/geppetto && make install"
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
