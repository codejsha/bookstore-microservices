import "@typespec/http";

using Http;

namespace SupportService;

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
  ): FaqFindAllResponse | UnprocessableEntityError;

  @doc("Create FAQ")
  @post
  create(@body body: FaqCreateRequest): CreatedResponse | BadRequestError | UnprocessableEntityError;

  @doc("Find FAQ")
  @route("{uid}")
  @get
  read(
    @path uid: uuid,
    @query increment_view?: boolean,
  ): FaqFindResponse | NotFoundError | UnprocessableEntityError;

  @doc("Update FAQ")
  @route("{uid}")
  @patch(#{ implicitOptionality: false })
  update(
    @path uid: uuid,
    @body body: FaqUpdateRequest,
  ): FaqUpdateResponse | BadRequestError | NotFoundError | UnprocessableEntityError;

  @doc("Delete FAQ")
  @route("{uid}")
  @delete
  remove(@path uid: uuid): NoContentResponse | NotFoundError | UnprocessableEntityError;
}
