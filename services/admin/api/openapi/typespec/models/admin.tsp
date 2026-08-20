import "@typespec/http";

using Http;

namespace AdminService;

// ─── Identity ───────────────────────────────────────────────────────────────

@doc("The authenticated administrator, as projected by the mesh from the JWT")
model AdminIdentityResponse {
  @doc("JWT subject; equals the user's uid across services")
  uid: uuid;

  @format("email")
  email?: string;

  name?: string;

  @doc("Realm roles carried by the token")
  roles: string[];
}

// ─── Dashboard ──────────────────────────────────────────────────────────────

@doc("A single downstream metric. `available` is false when that service could not be reached.")
model AdminMetric {
  @minValue(0)
  count: int64;

  available: boolean;
}

@doc("Cross-domain counts aggregated from the downstream services")
model AdminDashboardResponse {
  works: AdminMetric;
  orders: AdminMetric;
  warehouses: AdminMetric;
}
