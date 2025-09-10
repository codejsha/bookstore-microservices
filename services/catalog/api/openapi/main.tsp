import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/author.tsp";
import "./typespec/models/publisher.tsp";
import "./typespec/models/subject.tsp";
import "./typespec/models/work.tsp";
import "./typespec/models/edition.tsp";
import "./typespec/routes/author.tsp";
import "./typespec/routes/publisher.tsp";
import "./typespec/routes/subject.tsp";
import "./typespec/routes/work.tsp";
import "./typespec/routes/edition.tsp";

using Http;
using OpenAPI;
using Versioning;

@service(#{ title: "Catalog service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://catalog.bookstore.svc.cluster.local:8080", "catalog service")
@server("http://localhost:8080", "localhost")
@versioned(CatalogService.Versions)
namespace CatalogService;

enum Versions {
  v1_0_0: "1.0.0",
}
