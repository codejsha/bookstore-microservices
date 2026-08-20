import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/settlement.tsp";
import "./typespec/routes/settlement.tsp";

using Http;
using OpenAPI;
using Versioning;

/** Settlement service API of Bookstore microservices */
@service(#{ title: "Settlement service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://settlement.bookstore.svc.cluster.local:8080", "settlement service")
@server("http://localhost:8080", "localhost")
@versioned(SettlementService.Versions)
namespace SettlementService;

enum Versions {
  v1_0_0: "1.0.0",
}
