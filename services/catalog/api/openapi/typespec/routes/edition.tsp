import "@typespec/http";

using Http;

namespace CatalogService;

@route("/api/v1/editions")
@tag("Edition")
@useAuth(BearerAuth)
interface Editions {
  @doc("Search editions")
  @get
  @useAuth(NoAuth | BearerAuth)
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
  ): EditionFindAllResponse | BadRequestError;

  @doc("Create new edition")
  @post
  create(@body body: EditionCreateRequest): CreatedResponse | BadRequestError | NotFoundError | ConflictError | UnauthorizedError | ForbiddenError;

  @doc("Find edition with its work (uid and title only) and publisher")
  @route("{uid}")
  @get
  @useAuth(NoAuth | BearerAuth)
  read(@path uid: uuid): EditionFindResponse | NotFoundError;

  @doc("Update edition")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: EditionUpdateRequest): EditionUpdateResponse | BadRequestError | NotFoundError | ConflictError | UnauthorizedError | ForbiddenError;
}
