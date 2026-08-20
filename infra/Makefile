SHELL   := bash
LIVE    := $(abspath terragrunt)
TG_ARGS := -input=false $(TG_FLAGS)
RUN_ALL := --non-interactive

BOOTSTRAP_L1 := service/calico
BOOTSTRAP_L2 := service/prometheus-crds service/cert-manager service/gateway-api
BOOTSTRAP_L3 := service/istio service/vault-transit
BOOTSTRAP_L4 := service/vault-transit-config
BOOTSTRAP_L5 := service/vault
BOOTSTRAP_L6 := service/seaweedfs
BOOTSTRAP    := $(BOOTSTRAP_L1) $(BOOTSTRAP_L2) $(BOOTSTRAP_L3) $(BOOTSTRAP_L4) $(BOOTSTRAP_L5) $(BOOTSTRAP_L6)

BOOTSTRAP_EXCLUDE := $(foreach m,$(BOOTSTRAP),--queue-exclude-dir $(LIVE)/$(m))

POST_RESOURCE         := service/oauth2-proxy
POST_RESOURCE_EXCLUDE := $(foreach m,$(POST_RESOURCE),--queue-exclude-dir $(LIVE)/$(m))

DISABLED         := service/youtrack
DISABLED_EXCLUDE := $(foreach m,$(DISABLED),--queue-exclude-dir $(LIVE)/$(m))

define bootstrap_layer
	@echo "==> bootstrap $(2) layer: $(1)"
	@pids=""; for mod in $(1); do \
	  ( cd "$(LIVE)/$$mod" && terragrunt $(2) -auto-approve $(TG_ARGS) ) & \
	  pids="$$pids $$!"; \
	done; \
	rc=0; for p in $$pids; do wait $$p || rc=1; done; \
	[ $$rc -eq 0 ] || { echo "==> bootstrap $(2) layer failed: $(1)"; exit 1; }
endef

.DEFAULT_GOAL := help

.PHONY: none
none:
	@echo "No target specified."

.PHONY: help
help:
	@echo "make apply [MODULE=...]            run-all apply (service/ then resource/ then post-resource), or one unit"
	@echo "make apply-service                 run-all apply service/ only (bootstrap + post-resource + disabled excluded)"
	@echo "make apply-resource                run-all apply resource/ only"
	@echo "make apply-post-resource           apply post-resource units only (oauth2-proxy)"
	@echo "make apply-bootstrap-base          bootstrap L1-L3 (calico/crds/istio/vault-transit helm) — no unsealer access needed"
	@echo "make apply-bootstrap-rest          bootstrap L4-L6 (vault-transit-config/vault/seaweedfs) — needs unsealer port-forward + VAULT_UNSEALER_TOKEN"
	@echo "make apply-bootstrap               base + rest in one go (unsealer must already be reachable)"
	@echo "make apply-all                     apply-bootstrap + apply"
	@echo "make destroy [MODULE=...]          run-all destroy (post-resource then resource/ then service/), or one unit"
	@echo "make destroy-bootstrap             destroy bootstrap units (reverse order)"
	@echo "make destroy-all                   destroy + destroy-bootstrap"

.PHONY: apply-bootstrap-base
apply-bootstrap-base:
	$(call bootstrap_layer,$(BOOTSTRAP_L1),apply)
	$(call bootstrap_layer,$(BOOTSTRAP_L2),apply)
	$(call bootstrap_layer,$(BOOTSTRAP_L3),apply)

.PHONY: check-unsealer
check-unsealer:
	@[ -n "$$VAULT_UNSEALER_TOKEN" ] || { \
	  echo "VAULT_UNSEALER_TOKEN is not set."; \
	  echo "  export VAULT_UNSEALER_TOKEN=\$$(jq -r .root_token service/vault-transit/unsealer-keys.json)"; \
	  exit 1; }
	@curl -fsS -m 3 "$${VAULT_UNSEALER_ADDR:-http://127.0.0.1:8200}/v1/sys/health?standbyok=true" >/dev/null || { \
	  echo "Unsealer Vault unreachable, sealed, or uninitialized at $${VAULT_UNSEALER_ADDR:-http://127.0.0.1:8200}."; \
	  echo "  kubectl port-forward -n vault-transit vault-0 8200:8200"; \
	  echo "  (fresh install: vault operator init > service/vault-transit/unsealer-keys.json; then service/vault-transit/unseal.sh)"; \
	  exit 1; }

.PHONY: apply-bootstrap-rest
apply-bootstrap-rest: check-unsealer
	$(call bootstrap_layer,$(BOOTSTRAP_L4),apply)
	$(call bootstrap_layer,$(BOOTSTRAP_L5),apply)
	$(call bootstrap_layer,$(BOOTSTRAP_L6),apply)

.PHONY: apply-bootstrap
apply-bootstrap: apply-bootstrap-base apply-bootstrap-rest

.PHONY: destroy-bootstrap-rest
destroy-bootstrap-rest: check-unsealer
	$(call bootstrap_layer,$(BOOTSTRAP_L6),destroy)
	$(call bootstrap_layer,$(BOOTSTRAP_L5),destroy)
	$(call bootstrap_layer,$(BOOTSTRAP_L4),destroy)

.PHONY: destroy-bootstrap-base
destroy-bootstrap-base:
	$(call bootstrap_layer,$(BOOTSTRAP_L3),destroy)
	$(call bootstrap_layer,$(BOOTSTRAP_L2),destroy)
	$(call bootstrap_layer,$(BOOTSTRAP_L1),destroy)

.PHONY: destroy-bootstrap
destroy-bootstrap: destroy-bootstrap-rest destroy-bootstrap-base

define tree_run_all
	@echo "==> run-all $(2) $(1)/"
	@( cd "$(LIVE)/$(1)" && terragrunt run --all $(RUN_ALL) $(BOOTSTRAP_EXCLUDE) $(POST_RESOURCE_EXCLUDE) $(DISABLED_EXCLUDE) -- $(2) -auto-approve $(TG_ARGS) )
endef

.PHONY: apply-service
apply-service:
	$(call tree_run_all,service,apply)

.PHONY: apply-resource
apply-resource:
	$(call tree_run_all,resource,apply)

.PHONY: apply-post-resource
apply-post-resource:
	@for mod in $(POST_RESOURCE); do \
	  echo "==> apply post-resource: $$mod"; \
	  ( cd "$(LIVE)/$$mod" && terragrunt apply -auto-approve $(TG_ARGS) ); \
	done

.PHONY: apply
apply:
ifdef MODULE
	@cd "$(LIVE)/$(MODULE)" && terragrunt apply -auto-approve $(TG_ARGS)
else
	@$(MAKE) apply-service apply-resource apply-post-resource
endif

.PHONY: destroy-post-resource
destroy-post-resource:
	@for mod in $(POST_RESOURCE); do \
	  echo "==> destroy post-resource: $$mod"; \
	  ( cd "$(LIVE)/$$mod" && terragrunt destroy -auto-approve $(TG_ARGS) ); \
	done

.PHONY: destroy-resource
destroy-resource:
	$(call tree_run_all,resource,destroy)

.PHONY: destroy-service
destroy-service:
	$(call tree_run_all,service,destroy)

.PHONY: destroy
destroy:
ifdef MODULE
	@cd "$(LIVE)/$(MODULE)" && terragrunt destroy -auto-approve $(TG_ARGS)
else
	@$(MAKE) destroy-post-resource destroy-resource destroy-service
endif

.PHONY: apply-all
apply-all: apply-bootstrap apply
.PHONY: destroy-all
destroy-all: destroy destroy-bootstrap
