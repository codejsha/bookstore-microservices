import "@typespec/http";

using Http;

namespace IdentityService;

@route("/api/v1/risk")
@tag("Risk")
@useAuth(BearerAuth)
interface Risk {
  @doc("List all currently flagged principals")
  @get
  getAll(): RiskEntryListResponse;

  @doc("Flag a principal (restrict = deny writes, block = deny all)")
  @route("{uid}")
  @put
  flagPrincipal(@path uid: uuid, @body body: RiskFlagRequest): RiskEntryResponse | BadRequestError;

  @doc("Remove a principal's risk flag")
  @route("{uid}")
  @delete
  unflagPrincipal(@path uid: uuid): {
    @statusCode statusCode: 204;
  };
}
