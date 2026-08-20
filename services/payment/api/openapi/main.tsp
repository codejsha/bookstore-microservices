import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/customer.tsp";
import "./typespec/models/payment-method.tsp";
import "./typespec/models/payment.tsp";
import "./typespec/models/mandate.tsp";
import "./typespec/models/refund.tsp";
import "./typespec/routes/customer.tsp";
import "./typespec/routes/payment.tsp";
import "./typespec/routes/mandate.tsp";
import "./typespec/routes/refund.tsp";

using Http;
using OpenAPI;
using Versioning;

/** Payment service API of Bookstore microservices */
@service(#{ title: "Payment service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://payment.bookstore.svc.cluster.local:8080", "payment service")
@server("http://localhost:8080", "localhost")
@versioned(PaymentService.Versions)
namespace PaymentService;

enum Versions {
  v1_0_0: "1.0.0",
}
