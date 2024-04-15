test: generate manifests golint helm-test unit-test

manifests: controller-gen
	$(CONTROLLER_GEN) crd rbac:roleName=manager-role webhook paths="./.../..." output:crd:artifacts:config=config/crd/bases
	cp config/crd/bases/config.haproxy.com*.yaml helm/haproxy-operator/crds/

generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack\\boilerplate.go.txt" paths="./.../..."

golint: colanci-lint-bin
	$(GOLANGCI_LINT) run

unit-test: ginkgo-bin
	$(GINKGO) --no-color -r --randomize-all --randomize-suites --nodes=4 --compilers=4 --vet off

helm-test:
	helm lint helm/haproxy-operator
	helm template chart helm/haproxy-operator > /dev/null

CONTROLLER_GEN = bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.12.0)

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
	$(call go-get-tool,$(GINKGO),github.com/onsi/ginkgo/v2/ginkgo@v2.15.0)
