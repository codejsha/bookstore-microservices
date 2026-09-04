import "@typespec/http";

using Http;

namespace SupportService;

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
  ): TicketFindAllResponse | UnprocessableEntityError;

  @doc("Open new ticket")
  @post
  create(@body body: TicketCreateRequest): CreatedResponse | BadRequestError | UnprocessableEntityError;

  @doc("Find ticket")
  @route("{uid}")
  @get
  read(@path uid: uuid): TicketFindResponse | NotFoundError | UnprocessableEntityError;

  @doc("Update ticket")
  @route("{uid}")
  @patch(#{ implicitOptionality: false })
  update(
    @path uid: uuid,
    @body body: TicketUpdateRequest,
  ): TicketUpdateResponse | BadRequestError | NotFoundError | UnprocessableEntityError;

  @doc("Update ticket status")
  @route("{uid}/status")
  @patch(#{ implicitOptionality: false })
  updateStatus(
    @path uid: uuid,
    @body body: TicketUpdateStatusRequest,
  ): TicketUpdateResponse | NotFoundError | UnprocessableEntityError;
}
