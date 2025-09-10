import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/shipment.tsp";
import "./typespec/models/carrier.tsp";
import "./typespec/models/freight.tsp";
import "./typespec/models/stats.tsp";
import "./typespec/routes/shipment.tsp";
import "./typespec/routes/carrier.tsp";
import "./typespec/routes/freight.tsp";
import "./typespec/routes/stats.tsp";

using Http;
using OpenAPI;
using Versioning;

/** Delivery service API of Bookstore microservices */
@service(#{ title: "Delivery service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://delivery.bookstore.svc.cluster.local:8080", "delivery service")
@server("http://localhost:8080", "localhost")
@versioned(DeliveryService.Versions)
namespace DeliveryService;

enum Versions {
  v1_0_0: "1.0.0",
}
