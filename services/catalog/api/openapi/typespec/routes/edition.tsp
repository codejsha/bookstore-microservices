import "@typespec/http";

using Http;

namespace CatalogService;

@route("/api/v1/editions")
@tag("Edition")
@useAuth(BearerAuth)
interface Editions {
  @doc("Search editions")
  @get
  search(
    @query title?: string,
    @query isbn?: string,
    @query work_uid?: uuid,
    @query publisher_uid?: uuid,
    @query language?: string,
    @query ol_key?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): EditionFindAllResponse;

  @doc("Create new edition")
  @post
  create(@body body: EditionCreateRequest): CreatedResponse | BadRequestError | NotFoundError | ConflictError;

  @doc("Find edition with full details")
  @route("{uid}")
  @get
  read(@path uid: uuid): EditionFindResponse | NotFoundError;

  @doc("Update edition")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: EditionUpdateRequest): EditionUpdateResponse | BadRequestError | NotFoundError | ConflictError;
}
