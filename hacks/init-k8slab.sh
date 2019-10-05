#!/bin/bash
#
# initialize k8slab container image

set -ex

# set up k8slab user, paths, aliases
function setupUser {
  groupadd -g 7777 k8slab
  useradd -m -u 7777 -g k8slab k8slab
  chown k8slab:k8slab /home/k8slab

  cat << EOF >> /home/k8slab/.bashrc
export KUBECONFIG=/home/k8slab/.kube/kind-config-k8s-lab
alias k='kubectl'
EOF

}

# install stern
function installStern {
  curl -s -L -o /tmp/stern "https://github.com/wercker/stern/releases/download/${STERN_VERSION}/stern_linux_amd64"
  mv /tmp/stern /usr/local/bin
  chmod 755 /usr/local/bin/stern
}

# install krew
installKrew() {
  (
    set -x; cd "$(mktemp -d --tmpdir=/work)" &&
    curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/download/v0.3.1/krew.{tar.gz,yaml}" &&
    tar zxvf krew.tar.gz &&
    su k8slab -c "./krew-\"$(uname | tr '[:upper:]' '[:lower:]')_amd64\" install --manifest=krew.yaml --archive=krew.tar.gz"
  )
  cat << EOF >> /home/k8slab/.bashrc
export PATH="${KREW_ROOT:-/home/k8slab/.krew}/bin:$PATH"
EOF
}
export -f installKrew


setupUser
installStern
installKrew


exit 0
