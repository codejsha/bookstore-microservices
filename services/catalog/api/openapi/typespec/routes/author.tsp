import "@typespec/http";

using Http;

namespace CatalogService;

@route("/api/v1/authors")
@tag("Author")
@useAuth(BearerAuth)
interface Authors {
  @doc("Find all authors")
  @get
  getAll(
    @query name?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AuthorFindAllResponse;

  @doc("Register new author")
  @post
  create(@body body: AuthorCreateRequest): CreatedResponse;

  @doc("Find author")
  @route("{uid}")
  @get
  read(@path uid: uuid): AuthorFindResponse | NotFoundError;

  @doc("Update author")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: AuthorUpdateRequest): AuthorUpdateResponse | NotFoundError;
}
