SERVICE_NAME := admin-web
VERSION := 1.0.0

.PHONY: none
none:

.PHONY: version
version:
	@echo $(VERSION)

NEXUS_NPM_AUTH ?=
.PHONY: podman-build
podman-build:
	@auth=$$(mktemp); trap 'rm -f "$$auth" ca-chain.pem' EXIT; \
	{ if [ -n "$(NEXUS_NPM_AUTH)" ]; then printf '%s' "$(NEXUS_NPM_AUTH)"; \
	  else printf '%s:%s' "$$(vault kv get -mount=kv -field=username nexus/publisher)" \
	                      "$$(vault kv get -mount=kv -field=password nexus/publisher)" | base64; fi; } > "$$auth" 2>/dev/null; \
	[ -s "$$auth" ] || { echo "ERROR: no Nexus auth (set NEXUS_NPM_AUTH or ensure Vault kv/nexus/publisher is reachable)" >&2; exit 1; }; \
	vault read -field=certificate pki_int/cert/ca_chain > ca-chain.pem 2>/dev/null; \
	[ -s ca-chain.pem ] || { echo "ERROR: no Nexus CA (ensure Vault pki_int is reachable)" >&2; exit 1; }; \
	podman build \
		--no-cache \
		--platform linux/amd64 \
		--file ./ci/Dockerfile \
		--secret id=nexus_npm_auth,src="$$auth" \
		--tag harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION) .

.PHONY: podman-push
podman-push: podman-build
	podman push harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION)
