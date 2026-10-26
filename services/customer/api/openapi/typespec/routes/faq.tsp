import "@typespec/http";

using Http;

namespace CustomerService;

// ─── FAQs ───────────────────────────────────────────────────────────────────

@route("/api/v1/faqs")
@tag("Faq")
@useAuth(BearerAuth)
interface Faqs {
  @doc("Search FAQs")
  @get
  getAll(
    @query category_uid?: uuid,
    @query q?: string,
    @query published_only?: boolean,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): FaqFindAllResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Create a FAQ")
  @post
  create(@body body: FaqCreateRequest):
    | CreatedResponse
    | BadRequestError
    | UnauthorizedError
    | ForbiddenError;

  @doc("Find a FAQ")
  @route("{uid}")
  @get
  read(
    @path uid: uuid,
    @query increment_view?: boolean,
  ): FaqFindResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Update a FAQ")
  @route("{uid}")
  @patch(#{ implicitOptionality: false })
  update(
    @path uid: uuid,
    @body body: FaqUpdateRequest,
  ):
    | FaqFindResponse
    | BadRequestError
    | NotFoundError
    | UnauthorizedError
    | ForbiddenError;

  @doc("Delete a FAQ")
  @route("{uid}")
  @delete
  remove(@path uid: uuid): NoContentResponse | NotFoundError | UnauthorizedError | ForbiddenError;
}
