import "@typespec/http";

using Http;

namespace CatalogService;

@route("/api/v1/works")
@tag("Work")
@useAuth(BearerAuth)
interface Works {
  @doc("Search works. Served from the OpenSearch search index, not the primary database: a newly created or updated work appears here only after the CDC pipeline has indexed it")
  @get
  @useAuth(NoAuth | BearerAuth)
  search(
    @query title?: string,
    @query author_uid?: uuid,
    @query subject_uid?: uuid,
    @query ol_key?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): WorkFindAllResponse | BadRequestError;

  @doc("Create new work")
  @post
  create(@body body: WorkCreateRequest): CreatedResponse | BadRequestError | NotFoundError | ConflictError | UnauthorizedError | ForbiddenError;

  @doc("Find work with full details")
  @route("{uid}")
  @get
  @useAuth(NoAuth | BearerAuth)
  read(@path uid: uuid): WorkFindResponse | NotFoundError;

  @doc("Update work")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: WorkUpdateRequest): WorkUpdateResponse | BadRequestError | NotFoundError | ConflictError | UnauthorizedError | ForbiddenError;

  @doc("Find editions of a work")
  @route("{uid}/editions")
  @get
  @useAuth(NoAuth | BearerAuth)
  editions(
    @path uid: uuid,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): EditionFindAllResponse | BadRequestError;

  @doc("Full-text search backed by OpenSearch")
  @route("/search")
  @get
  @useAuth(NoAuth | BearerAuth)
  fullTextSearch(
    @query q?: string,
    @query author_uid?: uuid,
    @query subject_uid?: uuid,
    @query size?: Size,
    @query page?: Page,
  ): WorkSearchResponse | BadRequestError;
}
