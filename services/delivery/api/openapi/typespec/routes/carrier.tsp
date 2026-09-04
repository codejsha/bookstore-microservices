import "@typespec/http";

using Http;

namespace DeliveryService;

@route("/api/v1/carriers")
@tag("Carrier")
@useAuth(BearerAuth)
interface Carriers {
  @doc("List carriers")
  @get
  getAll(
    @query name?: string,
    @query status?: CarrierStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): CarrierFindAllResponse | UnprocessableEntityError;

  @doc("Register new carrier")
  @post
  create(@body body: CarrierCreateRequest): CreatedResponse | BadRequestError | ConflictError | UnprocessableEntityError;

  @doc("Find carrier")
  @route("{uid}")
  @get
  read(@path uid: uuid): CarrierFindResponse | NotFoundError | UnprocessableEntityError;

  @doc("Update carrier")
  @route("{uid}")
  @put
  update(
    @path uid: uuid,
    @body body: CarrierUpdateRequest,
  ): CarrierUpdateResponse | NotFoundError | BadRequestError | UnprocessableEntityError;
}
