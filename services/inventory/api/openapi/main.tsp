import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/stock.tsp";
import "./typespec/models/warehouse.tsp";
import "./typespec/models/transfer.tsp";
import "./typespec/models/audit.tsp";
import "./typespec/models/closing.tsp";
import "./typespec/models/balance.tsp";
import "./typespec/routes/stock.tsp";
import "./typespec/routes/warehouse.tsp";
import "./typespec/routes/transfer.tsp";
import "./typespec/routes/audit.tsp";
import "./typespec/routes/closing.tsp";
import "./typespec/routes/balance.tsp";

using Http;
using OpenAPI;
using Versioning;

/** Inventory service API of Bookstore microservices */
@service(#{ title: "Inventory service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://inventory.bookstore.svc.cluster.local:8080", "inventory service")
@server("http://localhost:8080", "localhost")
@versioned(InventoryService.Versions)
namespace InventoryService;

enum Versions {
  v1_0_0: "1.0.0",
}
