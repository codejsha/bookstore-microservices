import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/risk.tsp";
import "./typespec/models/user.tsp";
import "./typespec/routes/risk.tsp";
import "./typespec/routes/user.tsp";

using Http;
using OpenAPI;
using Versioning;

/** Identity service API of Bookstore microservices */
@service(#{ title: "Identity service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://identity.bookstore.svc.cluster.local:8080", "identity service")
@server("http://localhost:8080", "localhost")
@versioned(IdentityService.Versions)
namespace IdentityService;

enum Versions {
  v1_0_0: "1.0.0",
}
