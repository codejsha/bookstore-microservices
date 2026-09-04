import "@typespec/http";

using Http;

namespace NotificationService;

@route("/api/v1/templates")
@tag("Template")
@useAuth(BearerAuth)
interface Templates {
  @doc("Find all templates")
  @get
  getAll(
    @query notification_type?: NotificationType,
    @query channel?: Channel,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): TemplateFindAllResponse | UnprocessableEntityError;

  @doc("Create new template")
  @post
  create(@body body: TemplateCreateRequest): CreatedResponse | ForbiddenError | ConflictError | UnprocessableEntityError;

  @doc("Find template")
  @route("{uid}")
  @get
  read(@path uid: uuid): TemplateFindResponse | NotFoundError | UnprocessableEntityError;

  @doc("Update template")
  @route("{uid}")
  @put
  update(
    @path uid: uuid,
    @body body: TemplateUpdateRequest,
  ): TemplateUpdateResponse | ForbiddenError | NotFoundError | ConflictError | UnprocessableEntityError;

  @doc("Delete template")
  @route("{uid}")
  @delete
  delete(@path uid: uuid): NoContentResponse | ForbiddenError | NotFoundError | UnprocessableEntityError;
}
