import "@typespec/http";

using Http;

namespace AdminService;

@route("/api/v1/admin")
@tag("Admin")
@useAuth(BearerAuth)
interface Admin {
  @doc("The authenticated administrator. The admin console has no token of its own to inspect, so this is its only source of identity and roles.")
  @route("me")
  @get
  me(): AdminIdentityResponse | UnauthorizedError | ForbiddenError;

  @doc("Cross-domain counts for the console landing page")
  @route("dashboard")
  @get
  dashboard(): AdminDashboardResponse | UnauthorizedError | ForbiddenError;
}
