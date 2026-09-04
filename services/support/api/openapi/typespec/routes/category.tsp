import "@typespec/http";

using Http;

namespace SupportService;

@route("/api/v1/categories")
@tag("Category")
@useAuth(BearerAuth)
interface Categories {
  @doc("List categories")
  @get
  getAll(): CategoryFindAllResponse | UnprocessableEntityError;

  @doc("Create category")
  @post
  create(@body body: CategoryCreateRequest): CreatedResponse | BadRequestError | ConflictError | UnprocessableEntityError;

  @doc("Find category")
  @route("{uid}")
  @get
  read(@path uid: uuid): CategoryFindResponse | NotFoundError | UnprocessableEntityError;

  @doc("Update category")
  @route("{uid}")
  @patch(#{ implicitOptionality: false })
  update(
    @path uid: uuid,
    @body body: CategoryUpdateRequest,
  ): CategoryUpdateResponse | BadRequestError | NotFoundError | ConflictError | UnprocessableEntityError;

  @doc("Delete category")
  @route("{uid}")
  @delete
  remove(@path uid: uuid): NoContentResponse | NotFoundError | UnprocessableEntityError;
}
