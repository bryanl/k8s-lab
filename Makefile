.PHONY: install
install:
	@go install github.com/bryanl/k8s-lab/cmd/k8s-lab

.PHONY: build-wrapper
build-wrapper:
	@docker build -t bryanl/k8s-lab-wrapper -f Dockerfile.k8s-lab .

.PHONY: push-wrapper
push-wrapper:
	@docker push bryanl/k8s-lab-wrapper

.PHONY: build-lab-shell
build-lab-shell:
	@docker build -t bryanl/k8s-lab .

.PHONY: push-lab-shell
push-lab-shell:
	@docker push bryanl/k8s-lab

.PHONY: build
build: build-wrapper build-lab-shell

.PHONY: push
push: push-wrapper push-lab-shell

.PHONY: release-local
release-local:
	 goreleaser --snapshot --skip-publish --rm-dist