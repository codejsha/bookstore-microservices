import "@typespec/http";

using Http;

namespace CatalogService;

@route("/api/v1/publishers")
@tag("Publisher")
@useAuth(BearerAuth)
interface Publishers {
  @doc("Find all publishers")
  @get
  @useAuth(NoAuth | BearerAuth)
  getAll(
    @query name?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): PublisherFindAllResponse | BadRequestError;

  @doc("Register new publisher")
  @post
  create(@body body: PublisherCreateRequest): CreatedResponse | BadRequestError | ConflictError | UnauthorizedError | ForbiddenError;

  @doc("Find publisher")
  @route("{uid}")
  @get
  @useAuth(NoAuth | BearerAuth)
  read(@path uid: uuid): PublisherFindResponse | NotFoundError;

  @doc("Update publisher")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: PublisherUpdateRequest): PublisherUpdateResponse | BadRequestError | NotFoundError | ConflictError | UnauthorizedError | ForbiddenError;
}
