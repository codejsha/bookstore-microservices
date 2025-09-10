import "@typespec/http";

using Http;

namespace SupportService;

@route("/api/v1/categories")
@tag("Category")
@useAuth(BearerAuth)
interface Categories {
  @doc("List categories")
  @get
  getAll(): CategoryFindAllResponse;

  @doc("Create category")
  @post
  create(@body body: CategoryCreateRequest): CreatedResponse;

  @doc("Find category")
  @route("{uid}")
  @get
  read(@path uid: uuid): CategoryFindResponse | NotFoundError;

  @doc("Update category")
  @route("{uid}")
  @patch(#{ implicitOptionality: false })
  update(
    @path uid: uuid,
    @body body: CategoryUpdateRequest,
  ): CategoryUpdateResponse | NotFoundError;

  @doc("Delete category")
  @route("{uid}")
  @delete
  remove(@path uid: uuid): NoContentResponse | NotFoundError;
}
