import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/customer.tsp";
import "./typespec/models/wishlist.tsp";
import "./typespec/models/point.tsp";
import "./typespec/models/order.tsp";
import "./typespec/models/payment.tsp";
import "./typespec/models/delivery.tsp";
import "./typespec/routes/customer.tsp";
import "./typespec/routes/wishlist.tsp";
import "./typespec/routes/point.tsp";
import "./typespec/routes/order.tsp";
import "./typespec/routes/payment.tsp";
import "./typespec/routes/delivery.tsp";
import "./typespec/models/review.tsp";
import "./typespec/routes/review.tsp";

using Http;
using OpenAPI;
using Versioning;

/** Customer service API of Bookstore microservices */
@service(#{ title: "Customer service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://customer.bookstore.svc.cluster.local:8080", "customer service")
@server("http://localhost:8080", "localhost")
@versioned(CustomerService.Versions)
namespace CustomerService;

enum Versions {
  v1_0_0: "1.0.0",
}
