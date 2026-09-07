ARG GO_VERSION=1.26

ARG PROTOC_VERSION=35.1
ARG PROTOC_SHA256=6930ebf62bd4ea607b98fff052596c6ee564b9835b4ce172c75a3f53ae9d91b7
ARG BUF_VERSION=1.72.0
ARG BUF_SHA256=8720830e26a733da55bb89bcd3cb44849c0965fc0c44fb5d691cccdc64dca5af
ARG GRPC_JAVA_VERSION=1.78.0
ARG GRPC_JAVA_SHA256=c02dd4d757b6794072a1675c2a9f232631b64c3d53c1ec19effb0504859a82ff
ARG PROTOC_GEN_GO_VERSION=v1.36.11
ARG PROTOC_GEN_GO_GRPC_VERSION=v1.6.2

FROM --platform=$BUILDPLATFORM docker.io/library/golang:${GO_VERSION}-alpine AS goplugins

ARG PROTOC_GEN_GO_VERSION
ARG PROTOC_GEN_GO_GRPC_VERSION

WORKDIR /build

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod init bookstore/protoc-plugins \
    && go get google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION} \
    && go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION} \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o /out/protoc-gen-go google.golang.org/protobuf/cmd/protoc-gen-go \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o /out/protoc-gen-go-grpc google.golang.org/grpc/cmd/protoc-gen-go-grpc

FROM --platform=$BUILDPLATFORM docker.io/library/debian:bookworm-slim AS fetch

ARG PROTOC_VERSION
ARG PROTOC_SHA256
ARG BUF_VERSION
ARG BUF_SHA256
ARG GRPC_JAVA_VERSION
ARG GRPC_JAVA_SHA256

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl unzip \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /dl

RUN curl -fsSL -o protoc.zip "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" \
    && echo "${PROTOC_SHA256}  protoc.zip" | sha256sum -c - \
    && unzip -q protoc.zip -d protoc \
    && curl -fsSL -o buf "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-Linux-x86_64" \
    && echo "${BUF_SHA256}  buf" | sha256sum -c - \
    && curl -fsSL -o protoc-gen-grpc-java "https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/${GRPC_JAVA_VERSION}/protoc-gen-grpc-java-${GRPC_JAVA_VERSION}-linux-x86_64.exe" \
    && echo "${GRPC_JAVA_SHA256}  protoc-gen-grpc-java" | sha256sum -c - \
    && chmod 0755 protoc/bin/protoc buf protoc-gen-grpc-java

FROM docker.io/library/debian:bookworm-slim

ARG PROTOC_VERSION
ARG BUF_VERSION
ARG GRPC_JAVA_VERSION
ARG PROTOC_GEN_GO_VERSION
ARG PROTOC_GEN_GO_GRPC_VERSION

LABEL org.opencontainers.image.title="protoc-toolchain" \
      org.opencontainers.image.description="buf + protoc + protoc-gen-go/go-grpc/grpc-java pinned to the versions the committed gRPC stubs were generated with" \
      org.opencontainers.image.source="https://github.com/codejsha/bookstore-microservices" \
      dev.bookstore.protoc.version="${PROTOC_VERSION}" \
      dev.bookstore.buf.version="${BUF_VERSION}" \
      dev.bookstore.protoc-gen-go.version="${PROTOC_GEN_GO_VERSION}" \
      dev.bookstore.protoc-gen-go-grpc.version="${PROTOC_GEN_GO_GRPC_VERSION}" \
      dev.bookstore.protoc-gen-grpc-java.version="${GRPC_JAVA_VERSION}"

COPY --from=fetch /dl/protoc/bin/protoc /usr/local/bin/protoc
COPY --from=fetch /dl/protoc/include /usr/local/include
COPY --from=fetch /dl/buf /usr/local/bin/buf
COPY --from=fetch /dl/protoc-gen-grpc-java /usr/local/bin/protoc-gen-grpc-java
COPY --from=goplugins /out/protoc-gen-go /usr/local/bin/protoc-gen-go
COPY --from=goplugins /out/protoc-gen-go-grpc /usr/local/bin/protoc-gen-go-grpc

ENV BUF_CACHE_DIR=/tmp/buf-cache

WORKDIR /work

CMD ["buf", "--version"]
