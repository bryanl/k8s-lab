FROM fedora:30

RUN curl -s -o /usr/local/bin/kubectl -L https://storage.googleapis.com/kubernetes-release/release/v1.15.4/bin/linux/amd64/kubectl && \
    chmod 755 /usr/local/bin/kubectl

RUN groupadd -g 7777 k8slab && \
    useradd -r -u 7777 -g k8slab k8slab
USER k8slab

WORKDIR /home/k8slab