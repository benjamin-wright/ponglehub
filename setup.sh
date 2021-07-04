#!/usr/bin/bash

function install-direnv() {
    if cat ~/.bashrc | grep -q direnv; then
        echo "direnv already installed"
    else
        curl -sfL https://direnv.net/install.sh | bash
        echo 'eval "$(direnv hook bash)"' >> ~/.bashrc
        source ~/.bashrc
    fi
}

function install-earthly() {
    if [[ "$(which earthly)" == "" ]]; then
        sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/latest/download/earthly-linux-amd64 -O /usr/local/bin/earthly && chmod +x /usr/local/bin/earthly && /usr/local/bin/earthly bootstrap  --with-autocomplete'
    else
        echo "earthly already installed"
    fi
}

function install-k3d() {
    if [[ "$(which k3d)" == "" ]]; then
        curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
    else
        echo "k3d already installed"
    fi
}

function install-k9s() {
    if [[ "$(which k9s)" == "" ]]; then
        curl -sS https://webinstall.dev/k9s | bash
        echo 'export PATH="/home/bwright/.local/bin:$PATH"' >> ~/.bashrc
        source ~/.bashrc
    else
        echo "k9s already installed"
    fi
}

function install-helm() {
    if [[ "$(which helm)" == "" ]]; then
        curl -sS https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
    else
        echo "helm already installed"
    fi
}

function install-yq() {
    if [[ "$(which yq)" == "" ]]; then
        wget https://github.com/mikefarah/yq/releases/download/v4.9.6/yq_linux_amd64.tar.gz -O - |\
  tar xz && sudo mv yq_linux_amd64 /usr/bin/yq
    else
        echo "yq already installed"
    fi
}

function install-go() {
    if [[ "$(which go)" == "" ]]; then
        sudo apt-get update
        sudo apt install golang-go -y
    else
        echo "go already installed"
    fi
}

install-direnv
install-earthly
install-k3d
install-k9s
install-helm
install-yq
install-go