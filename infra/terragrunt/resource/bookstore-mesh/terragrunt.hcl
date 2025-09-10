include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

dependencies {
  paths = [
    "../bookstore-base",
    "../../service/istio",
    "../../service/gateway-api",
  ]
}

inputs = {
  namespace                    = "bookstore"
  ingress_gateway_principal    = "cluster.local/ns/istio-system/sa/bookstore-gateway-istio"
  hyperswitch_router_principal = "cluster.local/ns/hyperswitch/sa/hyperswitch-hyperswitch-router-role"
  keycloak_issuer              = "https://keycloak.example.com/realms/bookstore"
  keycloak_jwks_uri            = "http://keycloak-service.keycloak.svc.cluster.local:8080/realms/bookstore/protocol/openid-connect/certs"
  keycloak_audiences           = ["bookstore"]

  anonymous_paths = ["/health*", "/actuator/*", "/healthz", "/readyz", "/livez", "/oauth2/*"]

  anonymous_read_paths = [
    "/api/v1/works",
    "/api/v1/works/*",
    "/api/v1/editions",
    "/api/v1/editions/*",
    "/api/v1/authors",
    "/api/v1/authors/*",
    "/api/v1/publishers",
    "/api/v1/publishers/*",
    "/api/v1/subjects",
  ]

  request_timeout = "15s"
}

terraform {
  source = "${get_repo_root()}/infra//resource/bookstore-mesh"
}
