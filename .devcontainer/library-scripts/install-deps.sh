#!/usr/bin/env bash

set -o errexit -o pipefail

# Install k9s
curl -sSL https://github.com/derailed/k9s/releases/download/v0.24.13/k9s_Linux_x86_64.tar.gz -o k9s_Linux_x86_64.tar.gz
tar -xf k9s_Linux_x86_64.tar.gz k9s
sudo mv k9s /usr/local/bin/k9s
rm k9s_Linux_x86_64.tar.gz

# Install k3d
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash

# Install earthly
sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/latest/download/earthly-linux-amd64 -O /usr/local/bin/earthly && chmod +x /usr/local/bin/earthly'

# Install direnv
curl -sfL https://direnv.net/install.sh | bash
echo "eval \"\$(direnv hook bash)\"" >> ~/.bashrc

# Install helm
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash

# Install yq
wget https://github.com/mikefarah/yq/releases/download/v4.9.7/yq_linux_amd64.tar.gz -O - |\
  tar xz && sudo mv yq_linux_amd64 /usr/bin/yq