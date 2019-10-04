#!/bin/bash
#
# initialize k8slab container image

set -e

groupadd -g 7777 k8slab
useradd -m -u 7777 -g k8slab k8slab
chown k8slab:k8slab /home/k8slab

cat << EOF >> /home/k8slab/.bashrc
export KUBECONFIG=/home/k8slab/.kube/kind-config-k8s-lab
alias k='kubectl'
alias vim='nvim'
EOF

echo "k8slab ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/k8slab

curl -s -L -o /tmp/stern "https://github.com/wercker/stern/releases/download/${STERN_VERSION}/stern_linux_amd64"
mv /tmp/stern /usr/local/bin
chmod 755 /usr/local/bin/stern
