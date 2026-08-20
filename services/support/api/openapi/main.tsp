import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/ticket.tsp";
import "./typespec/models/comment.tsp";
import "./typespec/models/category.tsp";
import "./typespec/models/faq.tsp";
import "./typespec/routes/ticket.tsp";
import "./typespec/routes/comment.tsp";
import "./typespec/routes/category.tsp";
import "./typespec/routes/faq.tsp";

using Http;
using OpenAPI;
using Versioning;

/** Support service API of Bookstore microservices */
@service(#{ title: "Support service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://support.bookstore.svc.cluster.local:8080", "support service")
@server("http://localhost:8080", "localhost")
@versioned(SupportService.Versions)
namespace SupportService;

enum Versions {
  v1_0_0: "1.0.0",
}
