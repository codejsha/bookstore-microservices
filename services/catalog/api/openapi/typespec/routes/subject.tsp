import "@typespec/http";

using Http;

namespace CatalogService;

@route("/api/v1/subjects")
@tag("Subject")
@useAuth(BearerAuth)
interface Subjects {
  @doc("Find all subjects")
  @get
  @useAuth(NoAuth | BearerAuth)
  getAll(
    @query name?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): SubjectFindAllResponse | BadRequestError;

  @doc("Create subject")
  @post
  create(@body body: SubjectCreateRequest): CreatedResponse | BadRequestError | ConflictError | UnauthorizedError | ForbiddenError;

  @doc("Find subject")
  @route("{uid}")
  @get
  @useAuth(NoAuth | BearerAuth)
  read(@path uid: uuid): SubjectFindResponse | NotFoundError;

  @doc("Update subject")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: SubjectUpdateRequest): SubjectUpdateResponse | BadRequestError | NotFoundError | ConflictError | UnauthorizedError | ForbiddenError;
}
