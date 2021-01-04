variable "cluster" {
    type = string
}

variable "registry_port" {
    type = number
    default = 5000
}

variable "npm_port" {
    type = number
    default = 4873
}

variable "chartmuseum_port" {
    type = number
    default = 5002
}

variable "deploy_cluster" {
    type = bool
    default = true
}