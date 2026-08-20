import "@typespec/http";

using Http;

namespace CustomerService;

// ─── Customer aggregate root ───────────────────────────────────────────────

@route("/api/v1/customers")
@tag("Customer")
@useAuth(BearerAuth)
interface Customers {
  @doc("Find all customers")
  @get
  getAll(
    @query email?: string,
    @query name?: string,
    @query phone?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): CustomerFindAllResponse;

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
