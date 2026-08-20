LIVE_DIR := $(abspath ../../terragrunt/service/grafana)

.PHONY: none
none:
	@echo "No target specified."

.PHONY: init
init:
	cd $(LIVE_DIR) && terragrunt init

.PHONY: plan
plan: init
	cd $(LIVE_DIR) && terragrunt plan

.PHONY: apply
apply: init
	cd $(LIVE_DIR) && terragrunt apply -auto-approve

.PHONY: destroy
destroy:
	cd $(LIVE_DIR) && terragrunt destroy -auto-approve

.PHONY: clean
clean:
	rm -rf $(LIVE_DIR)/.terragrunt-cache
