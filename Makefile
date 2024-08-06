test: generate manifests docs golint helm-test unit-test

PATHS ?= "./..."
manifests: controller-gen
	$(CONTROLLER_GEN) crd rbac:roleName=manager-role webhook paths=${PATHS} output:crd:artifacts:config=config/crd/bases
	cp config/crd/bases/*.haproxy.com*.yaml helm/haproxy-operator/crds/

generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack\\boilerplate.go.txt" paths=${PATHS}

.PHONY: docs
docs:
	go run github.com/elastic/crd-ref-docs@v0.0.10 --config docs/config.yaml --renderer=markdown --output-path docs/api-reference.md

golint: colanci-lint-bin
	$(GOLANGCI_LINT) run

unit-test: ginkgo-bin
	$(GINKGO) --no-color -r --randomize-all --randomize-suites --nodes=4 --compilers=4 --vet off

helm-test:
	helm lint helm/haproxy-operator
	helm template chart helm/haproxy-operator > /dev/null

CONTROLLER_GEN = bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.15.0)

GOLANGCI_LINT = ./bin/golangci-lint
colanci-lint-bin:
	$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2)

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
}
endef

GINKGO = ./bin/ginkgo
ginkgo-bin:
	$(call go-get-tool,$(GINKGO),github.com/onsi/ginkgo/v2/ginkgo@v2.19.0)
