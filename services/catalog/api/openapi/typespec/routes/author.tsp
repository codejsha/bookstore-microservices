import "@typespec/http";

using Http;

namespace CatalogService;

@route("/api/v1/authors")
@tag("Author")
@useAuth(BearerAuth)
interface Authors {
  @doc("Find all authors")
  @get
  @useAuth(NoAuth | BearerAuth)
  getAll(
    @query name?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AuthorFindAllResponse | BadRequestError;

  @doc("Register new author")
  @post
  create(@body body: AuthorCreateRequest): CreatedResponse | BadRequestError | ConflictError | UnauthorizedError | ForbiddenError;

  @doc("Find author")
  @route("{uid}")
  @get
  @useAuth(NoAuth | BearerAuth)
  read(@path uid: uuid): AuthorFindResponse | NotFoundError;

  @doc("Update author")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: AuthorUpdateRequest): AuthorUpdateResponse | BadRequestError | NotFoundError | ConflictError | UnauthorizedError | ForbiddenError;
}
