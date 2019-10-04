FROM fedora:30

RUN curl -s -o /usr/local/bin/kubectl -L https://storage.googleapis.com/kubernetes-release/release/v1.15.4/bin/linux/amd64/kubectl && \
    chmod 755 /usr/local/bin/kubectl && \
    dnf install -q -y which neovim nano jq net-tools

ARG STERN_VERSION=1.11.0

COPY hacks/init-k8slab.sh /tmp/init-k8slab.sh
RUN chmod +x /tmp/init-k8slab.sh && \
    /tmp/init-k8slab.sh && \
    rm -rf /tmp/init-k8slab.sh

USER k8slab

COPY lab_scripts /home/k8slab/lab_scripts

WORKDIR /home/k8slab
