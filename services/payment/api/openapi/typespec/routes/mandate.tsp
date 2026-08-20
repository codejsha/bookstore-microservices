import "@typespec/http";

using Http;

namespace PaymentService;

@route("/api/v1/mandates")
@tag("Mandate")
@useAuth(BearerAuth)
interface Mandates {
  @doc("Find all mandates")
  @get
  getAll(
    @query customer_id?: string,
    @query mandate_status?: MandateStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): MandateFindAllResponse;

  @doc("Register a payment instrument and create the mandate that authorizes future charges")
  @post
  setup(@body body: MandateSetupRequest): MandateFindResponse | BadRequestError;

  @doc("Find mandate")
  @route("{uid}")
  @get
  read(@path uid: uuid): MandateFindResponse | NotFoundError;

  @doc("Revoke mandate")
  @route("{uid}/revoke")
  @post
  revoke(@path uid: uuid): MandateFindResponse | NotFoundError;
}
