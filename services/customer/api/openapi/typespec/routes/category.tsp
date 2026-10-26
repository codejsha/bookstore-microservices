import "@typespec/http";

using Http;

namespace CustomerService;

// ─── Ticket categories ──────────────────────────────────────────────────────

@route("/api/v1/categories")
@tag("Category")
@useAuth(BearerAuth)
interface Categories {
  @doc("List categories")
  @get
  getAll(): CategoryFindAllResponse | UnauthorizedError | ForbiddenError;

  @doc("Create a category")
  @post
  create(@body body: CategoryCreateRequest):
    | CreatedResponse
    | BadRequestError
    | UnauthorizedError
    | ForbiddenError
    | ConflictError;

  @doc("Find a category")
  @route("{uid}")
  @get
  read(@path uid: uuid): CategoryFindResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Update a category")
  @route("{uid}")
  @patch(#{ implicitOptionality: false })
  update(
    @path uid: uuid,
    @body body: CategoryUpdateRequest,
  ):
    | CategoryFindResponse
    | BadRequestError
    | NotFoundError
    | UnauthorizedError
    | ForbiddenError
    | ConflictError;

  @doc("Delete a category")
  @route("{uid}")
  @delete
  remove(@path uid: uuid): NoContentResponse | NotFoundError | UnauthorizedError | ForbiddenError;
}
