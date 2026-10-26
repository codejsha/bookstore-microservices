import "@typespec/http";

using Http;

namespace CustomerService;

// ─── Ticket comments ────────────────────────────────────────────────────────

@route("/api/v1/tickets/{uid}/comments")
@tag("Comment")
@useAuth(BearerAuth)
interface Comments {
  @doc("List ticket comments")
  @get
  getAll(
    @path uid: uuid,
    @query include_internal?: boolean,
  ): CommentFindAllResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Add a comment to a ticket")
  @post
  create(
    @path uid: uuid,
    @body body: CommentCreateRequest,
  ):
    | CreatedResponse
    | BadRequestError
    | NotFoundError
    | UnauthorizedError
    | ForbiddenError;
}
