import "@typespec/http";

using Http;

namespace AdminService;

@route("/api/v1/admin/risk")
@tag("Admin Risk")
@useAuth(BearerAuth)
interface AdminRisk {
  @doc("List all currently flagged principals")
  @get
  listRisk(): AdminRiskEntryListResponse | UnauthorizedError | ForbiddenError;

  @doc("Flag a principal (restrict = deny writes, block = deny all). An administrator may not flag themselves.")
  @route("{uid}")
  @put
  flagRisk(
    @path uid: uuid,
    @body body: AdminRiskFlagRequest,
  ): AdminRiskEntryResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Remove a principal's risk flag")
  @route("{uid}")
  @delete
  unflagRisk(@path uid: uuid): {
    @statusCode statusCode: 204;
  } | UnauthorizedError | ForbiddenError;
}
