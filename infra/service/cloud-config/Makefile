CHART_NAME := config-server
CHART_TREE := $(shell git rev-parse --short=8 HEAD:$(shell git rev-parse --show-prefix)helm)
CHART_VERSION := 1.0.0-$(CHART_TREE)
HELM_REGISTRY := oci://harbor.example.com/bookstore-helm-charts

.PHONY: none
none:
	@echo "No target specified."

.PHONY: version
version:
	@echo $(CHART_VERSION)

.PHONY: spec
spec:
	yq -i '.version = "$(CHART_VERSION)"' helm/Chart.yaml

.PHONY: helm-package
helm-package:
	helm package ./helm --version $(CHART_VERSION)

.PHONY: helm-push
helm-push: helm-package
	helm push $(CHART_NAME)-$(CHART_VERSION).tgz $(HELM_REGISTRY)

.PHONY: clean
clean:
	rm -f $(CHART_NAME)-*.tgz
