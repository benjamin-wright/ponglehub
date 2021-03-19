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

resource "null_resource" "npm" {
    depends_on = [
        docker_container.npm
    ]

    provisioner "local-exec" {
        command = "./scripts/setup-npm.sh"
    }

    provisioner "local-exec" {
        when = destroy
        command = "./scripts/restore-npm.sh"
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

resource "null_resource" "chartmuseum" {
    depends_on = [
        docker_container.chartmuseum
    ]

    provisioner "local-exec" {
        command = "./scripts/setup-chartmuseum.sh"
    }

    provisioner "local-exec" {
        when = destroy
        command = "./scripts/restore-chartmuseum.sh"
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