import "@typespec/http";

using Http;

namespace CatalogService;

@route("/api/v1/works")
@tag("Work")
@useAuth(BearerAuth)
interface Works {
  @doc("Search works")
  @get
  search(
    @query title?: string,
    @query author_uid?: uuid,
    @query subject_uid?: uuid,
    @query ol_key?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): WorkFindAllResponse;

  @doc("Create new work")
  @post
  create(@body body: WorkCreateRequest): CreatedResponse;

  @doc("Find work with full details")
  @route("{uid}")
  @get
  read(@path uid: uuid): WorkFindResponse | NotFoundError;

  @doc("Update work")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: WorkUpdateRequest): WorkUpdateResponse | NotFoundError;

  @doc("Find editions of a work")
  @route("{uid}/editions")
  @get
  editions(
    @path uid: uuid,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): EditionFindAllResponse;

  @doc("Full-text search backed by OpenSearch")
  @route("/search")
  @get
  fullTextSearch(
    @query q?: string,
    @query author_uid?: uuid,
    @query subject_uid?: uuid,
    @query size?: Size,
    @query page?: Page,
  ): WorkSearchResponse;
}
