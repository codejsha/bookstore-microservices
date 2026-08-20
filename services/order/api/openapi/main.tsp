import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/order.tsp";
import "./typespec/models/order-item.tsp";
import "./typespec/models/order-shipping.tsp";
import "./typespec/models/order-adjustment.tsp";
import "./typespec/models/cart.tsp";
import "./typespec/routes/order.tsp";
import "./typespec/routes/order-item.tsp";
import "./typespec/routes/order-shipping.tsp";
import "./typespec/routes/order-adjustment.tsp";
import "./typespec/routes/cart.tsp";

using Http;
using OpenAPI;
using Versioning;

/** Order service API of Bookstore microservices */
@service(#{ title: "Order service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://order.bookstore.svc.cluster.local:8080", "order service")
@server("http://localhost:8080", "localhost")
@versioned(OrderService.Versions)
namespace OrderService;

enum Versions {
  v1_0_0: "1.0.0",
}
