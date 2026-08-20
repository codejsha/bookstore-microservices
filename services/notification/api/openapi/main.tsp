import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";
import "@typespec/versioning";

import "./typespec/models/common.tsp";
import "./typespec/models/notification.tsp";
import "./typespec/models/template.tsp";
import "./typespec/routes/notification.tsp";
import "./typespec/routes/template.tsp";

using Http;
using OpenAPI;
using Versioning;

/** Notification service API of Bookstore microservices */
@service(#{ title: "Notification service API" })
@info(#{
  version: "{api-version}",
  contact: #{
    name: "Admin",
    email: "admin@example.com",
  },
})
@server("http://notification.bookstore.svc.cluster.local:8080", "notification service")
@server("http://localhost:8080", "localhost")
@versioned(NotificationService.Versions)
namespace NotificationService;

enum Versions {
  v1_0_0: "1.0.0",
}
