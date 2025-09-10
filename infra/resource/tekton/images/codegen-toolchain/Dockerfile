# ── In-house generators ─────────────────────────────────────────────────────
ARG CODEGEN_TOOLS_IMAGE=harbor.example.com/bookstore/codegen-tools:0.1.0
FROM ${CODEGEN_TOOLS_IMAGE} AS tools

# ── Base + OS tools ─────────────────────────────────────────────────────────
FROM docker.io/library/node:24-bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
      make git curl ca-certificates python3 \
    && rm -rf /var/lib/apt/lists/*

# ── pnpm (pinned to the repo's packageManager) ──────────────────────────────
RUN corepack enable && corepack prepare pnpm@10.33.0 --activate

# ── buf / redocly / yq ──────────────────────────────────────────────────────
ARG BUF_VERSION=1.71.0
ARG YQ_VERSION=4.53.3
ARG TARGETARCH=amd64
RUN curl -fsSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-Linux-x86_64" \
      -o /usr/local/bin/buf && chmod +x /usr/local/bin/buf \
 && curl -fsSL "https://github.com/mikefarah/yq/releases/download/v${YQ_VERSION}/yq_linux_${TARGETARCH}" \
      -o /usr/local/bin/yq && chmod +x /usr/local/bin/yq
ARG REDOCLY_VERSION=2.38.0
RUN npm install -g "@redocly/cli@${REDOCLY_VERSION}"

# ── In-house generators from the tools image ────────────────────────────────
COPY --from=tools /usr/local/bin/oapi-codegen-cli /usr/local/bin/oapi-codegen-cli
COPY --from=tools /usr/local/bin/gorm-codegen-cli /usr/local/bin/gorm-codegen-cli
RUN chmod +x /usr/local/bin/oapi-codegen-cli /usr/local/bin/gorm-codegen-cli

WORKDIR /workspace
