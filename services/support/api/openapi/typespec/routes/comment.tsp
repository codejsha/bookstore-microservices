import "@typespec/http";

using Http;

namespace SupportService;

@route("/api/v1/tickets/{ticket_uid}/comments")
@tag("Comment")
@useAuth(BearerAuth)
interface Comments {
  @doc("List ticket comments")
  @get
  getAll(
    @path ticket_uid: uuid,
    @query include_internal?: boolean,
  ): CommentFindAllResponse | NotFoundError | UnprocessableEntityError;

  @doc("Add comment to ticket")
  @post
  create(
    @path ticket_uid: uuid,
    @body body: CommentCreateRequest,
  ): CreatedResponse | NotFoundError | UnprocessableEntityError;
}
