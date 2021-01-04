locals {
    scratch_dir = "${abspath(path.module)}/.scratch"
}

data "docker_registry_image" "npm" {
    name = "verdaccio/verdaccio:4"
}

resource "docker_image" "npm" {
  name          = data.docker_registry_image.npm.name
  pull_triggers = [data.docker_registry_image.npm.sha256_digest]
}

resource "docker_container" "npm" {
    image = docker_image.npm.latest
    name = "${var.cluster}-npm"
    restart = "always"
    ports {
        internal = 4873
        external = var.npm_port
    }
}

data "docker_registry_image" "chartmuseum" {
    name = "chartmuseum/chartmuseum:v0.12.0"
}

resource "docker_image" "chartmuseum" {
  name          = data.docker_registry_image.chartmuseum.name
  pull_triggers = [data.docker_registry_image.chartmuseum.sha256_digest]
}

resource "docker_container" "chartmuseum" {
    image = docker_image.chartmuseum.latest
    name = "${var.cluster}-charts"
    restart = "always"
    ports {
        internal = 8080
        external = var.chartmuseum_port
    }
    env = [ "STORAGE=local", "STORAGE_LOCAL_ROOTDIR=/home/chartmuseum" ]
    volumes {
        volume_name = "${var.cluster}-charts"
        container_path = "/home/chartmuseum"
    }
}

data "docker_registry_image" "registry" {
  name = "registry:2"
}

resource "docker_image" "registry" {
  name          = data.docker_registry_image.registry.name
  pull_triggers = [data.docker_registry_image.registry.sha256_digest]
}

resource "docker_network" "registry" {
    name = "${var.cluster}-network"
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
    count = var.deploy_cluster ? 1 : 0

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
