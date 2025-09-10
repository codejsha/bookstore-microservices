import "@typespec/http";

using Http;

namespace AdminService;

@route("/api/v1/admin/catalog")
@tag("Admin Catalog")
@useAuth(BearerAuth)
interface AdminCatalog {
  @doc("List works")
  @route("works")
  @get
  listWorks(
    @query title?: string,
    @query author_uid?: uuid,
    @query subject_uid?: uuid,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AdminWorkFindAllResponse | UnauthorizedError | ForbiddenError;

  @doc("Create a work. The catalog answers 201 with no body, so the new uid is not returned.")
  @route("works")
  @post
  createWork(@body body: AdminWorkCreateRequest):
    | {
        @statusCode statusCode: 201;
      }
    | BadRequestError
    | UnauthorizedError
    | ForbiddenError;

  @doc("Find a work by uid")
  @route("works/{uid}")
  @get
  readWork(@path uid: uuid): AdminWorkResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Update a work")
  @route("works/{uid}")
  @put
  updateWork(
    @path uid: uuid,
    @body body: AdminWorkUpdateRequest,
  ): AdminWorkResponse | BadRequestError | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("List authors, for the work form's author picker")
  @route("authors")
  @get
  listAuthors(
    @query name?: string,
    @query size?: Size,
    @query page?: Page,
  ): AdminAuthorFindAllResponse | UnauthorizedError | ForbiddenError;

  @doc("List subjects, for the work form's subject picker")
  @route("subjects")
  @get
  listSubjects(
    @query name?: string,
    @query size?: Size,
    @query page?: Page,
  ): AdminSubjectFindAllResponse | UnauthorizedError | ForbiddenError;
}
