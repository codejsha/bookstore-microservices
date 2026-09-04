import "@typespec/http";

using Http;

namespace PaymentService;

// ─── Payment aggregate root ─────────────────────────────────────────────────

@route("/api/v1/payments")
@tag("Payment")
@useAuth(BearerAuth)
interface Payments {
  @doc("Find all payments")
  @get
  getAll(
    @query customer_id?: string,
    @query status?: PaymentStatus,
    @query connector?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): PaymentFindAllResponse;

  @doc("Create payment")
  @post
  create(@body body: PaymentCreateRequest): CreatedResponse | BadRequestError | ConflictError;

  @doc("Find payment")
  @route("{uid}")
  @get
  read(@path uid: uuid): PaymentFindResponse | NotFoundError;

  @doc("Update payment")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: PaymentUpdateRequest): PaymentUpdateResponse | BadRequestError | NotFoundError;
}

// ─── Payment attempts (Payment aggregate internal entity) ───────────────────

@route("/api/v1/payments/{payment_uid}/attempts")
@tag("Payment - Attempt")
@useAuth(BearerAuth)
interface PaymentAttempts {
  @doc("Find all attempts of a payment")
  @get
  getAll(
    @path payment_uid: uuid,
    @query status?: PaymentStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): PaymentAttemptFindAllResponse | NotFoundError;

  @doc("Find attempt of a payment")
  @route("{uid}")
  @get
  read(@path payment_uid: uuid, @path uid: uuid): PaymentAttemptFindResponse | NotFoundError;
}
