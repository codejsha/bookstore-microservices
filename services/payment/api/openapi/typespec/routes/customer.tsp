import "@typespec/http";

using Http;

namespace PaymentService;

// ─── Customer aggregate root ────────────────────────────────────────────────

@route("/api/v1/customers")
@tag("Customer")
@useAuth(BearerAuth)
interface Customers {
  @doc("Find all customers")
  @get
  getAll(
    @query customer_id?: string,
    @query email?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): CustomerFindAllResponse;

  @doc("Create customer")
  @post
  create(@body body: CustomerCreateRequest): CreatedResponse;

  @doc("Find customer")
  @route("{uid}")
  @get
  read(@path uid: uuid): CustomerFindResponse | NotFoundError;

  @doc("Update customer")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: CustomerUpdateRequest): CustomerUpdateResponse | NotFoundError;

  @doc("Delete customer")
  @route("{uid}")
  @delete
  delete(@path uid: uuid): NoContentResponse | NotFoundError;
}

// ─── Payment methods (Customer aggregate internal entity) ───────────────────

@route("/api/v1/customers/{customer_uid}/payment-methods")
@tag("Customer - Payment method")
@useAuth(BearerAuth)
interface CustomerPaymentMethods {
  @doc("Find all payment methods of a customer")
  @get
  getAll(
    @path customer_uid: uuid,
    @query payment_method?: PaymentMethodEnum,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): PaymentMethodFindAllResponse | NotFoundError;

  @doc("Create payment method for a customer")
  @post
  create(
    @path customer_uid: uuid,
    @body body: PaymentMethodCreateRequest,
  ): CreatedResponse | NotFoundError;

  @doc("Find payment method of a customer")
  @route("{uid}")
  @get
  read(@path customer_uid: uuid, @path uid: uuid): PaymentMethodFindResponse | NotFoundError;

  @doc("Update payment method of a customer")
  @route("{uid}")
  @put
  update(
    @path customer_uid: uuid,
    @path uid: uuid,
    @body body: PaymentMethodUpdateRequest,
  ): PaymentMethodUpdateResponse | NotFoundError;

  @doc("Delete payment method of a customer")
  @route("{uid}")
  @delete
  delete(@path customer_uid: uuid, @path uid: uuid): NoContentResponse | BadRequestError | NotFoundError;
}
