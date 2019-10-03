FROM fedora:30

RUN curl -s -o /usr/local/bin/kubectl -L https://storage.googleapis.com/kubernetes-release/release/v1.15.4/bin/linux/amd64/kubectl && \
    chmod 755 /usr/local/bin/kubectl && \
    dnf install -q -y which neovim nano jq net-tools

RUN groupadd -g 7777 k8slab && \
    useradd -m -u 7777 -g k8slab k8slab && \
    chown k8slab:k8slab /home/k8slab && \
    echo "export KUBECONFIG=/home/k8slab/.kube/kind-config-k8s-lab" >> /home/k8slab/.bashrc && \
    echo "alias k='kubectl'" >> /home/k8slab/.bashrc && \
    echo "alias vim='nvim'" >> /home/k8slab/.bashrc && \
    echo "k8slab ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/k8slab
USER k8slab

COPY lab_scripts /home/k8slab/lab_scripts

WORKDIR /home/k8slab
