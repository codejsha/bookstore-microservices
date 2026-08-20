import "@typespec/http";

using Http;

namespace CatalogService;

@route("/api/v1/publishers")
@tag("Publisher")
@useAuth(BearerAuth)
interface Publishers {
  @doc("Find all publishers")
  @get
  getAll(
    @query name?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): PublisherFindAllResponse;

  @doc("Register new publisher")
  @post
  create(@body body: PublisherCreateRequest): CreatedResponse;

  @doc("Find publisher")
  @route("{uid}")
  @get
  read(@path uid: uuid): PublisherFindResponse | NotFoundError;

  @doc("Update publisher")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: PublisherUpdateRequest): PublisherUpdateResponse | NotFoundError;
}
