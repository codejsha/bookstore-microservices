import "@typespec/http";

using Http;

namespace CustomerService;

// ─── Support tickets ────────────────────────────────────────────────────────

@route("/api/v1/tickets")
@tag("Ticket")
@useAuth(BearerAuth)
interface Tickets {
  @doc("List tickets")
  @get
  getAll(
    @query customer_uid?: uuid,
    @query assignee_uid?: uuid,
    @query category_uid?: uuid,
    @query status?: TicketStatus,
    @query priority?: TicketPriority,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): TicketFindAllResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Open a new ticket")
  @post
  create(@body body: TicketCreateRequest):
    | CreatedResponse
    | BadRequestError
    | UnauthorizedError
    | ForbiddenError;

  @doc("Find a ticket")
  @route("{uid}")
  @get
  read(@path uid: uuid): TicketFindResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Update a ticket")
  @route("{uid}")
  @patch(#{ implicitOptionality: false })
  update(
    @path uid: uuid,
    @body body: TicketUpdateRequest,
  ):
    | TicketFindResponse
    | BadRequestError
    | NotFoundError
    | UnauthorizedError
    | ForbiddenError;

  @doc("Update a ticket status")
  @route("{uid}/status")
  @patch(#{ implicitOptionality: false })
  updateStatus(
    @path uid: uuid,
    @body body: TicketUpdateStatusRequest,
  ):
    | TicketFindResponse
    | BadRequestError
    | NotFoundError
    | UnauthorizedError
    | ForbiddenError;
}
