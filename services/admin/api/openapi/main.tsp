import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/admin.tsp";
import "./typespec/models/catalog.tsp";
import "./typespec/models/inventory.tsp";
import "./typespec/models/order.tsp";
import "./typespec/models/payment.tsp";
import "./typespec/models/settlement.tsp";
import "./typespec/models/risk.tsp";
import "./typespec/models/user.tsp";
import "./typespec/routes/admin.tsp";
import "./typespec/routes/catalog.tsp";
import "./typespec/routes/inventory.tsp";
import "./typespec/routes/order.tsp";
import "./typespec/routes/payment.tsp";
import "./typespec/routes/settlement.tsp";
import "./typespec/routes/risk.tsp";
import "./typespec/routes/user.tsp";

using Http;
using OpenAPI;
using Versioning;

/** Admin BFF API of Bookstore microservices */
@service(#{ title: "Admin service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://admin.bookstore.svc.cluster.local:8080", "admin service")
@server("http://localhost:8080", "localhost")
@versioned(AdminService.Versions)
namespace AdminService;

enum Versions {
  v1_0_0: "1.0.0",
}
