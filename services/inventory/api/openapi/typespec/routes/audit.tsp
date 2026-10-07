import "@typespec/http";

using Http;

namespace InventoryService;

@route("/api/v1/audits")
@tag("Audit")
@useAuth(BearerAuth)
interface Audits {
  @doc("Find all stock audits")
  @get
  getAll(
    @query warehouse_uid?: uuid,
    @query status?: AuditStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AuditFindAllResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Create stock audit (실사)")
  @post
  create(@body body: AuditCreateRequest): AuditFindResponse | BadRequestError | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Find stock audit")
  @route("{uid}")
  @get
  read(@path uid: uuid): AuditFindResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Complete stock audit and apply adjustments")
  @route("{uid}/complete")
  @post
  complete(@path uid: uuid): AuditFindResponse | BadRequestError | NotFoundError | ConflictError | UnauthorizedError | ForbiddenError;
}
