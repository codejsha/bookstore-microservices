import "@typespec/http";

using Http;

namespace NotificationService;

@route("/api/v1/notifications")
@tag("Notification")
@useAuth(BearerAuth)
interface Notifications {
  @doc("Find all notifications")
  @get
  getAll(
    @query user_uid?: uuid,
    @query notification_type?: NotificationType,
    @query channel?: Channel,
    @query status?: NotificationStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): NotificationFindAllResponse | ForbiddenError;

  @doc("Send new notification")
  @post
  create(@body body: SendNotificationRequest): CreatedResponse | ForbiddenError;

  @doc("Find notification")
  @route("{uid}")
  @get
  read(@path uid: uuid): NotificationFindResponse | ForbiddenError | NotFoundError;

  @doc("Mark notification as sent")
  @route("{uid}/sent")
  @patch(#{ implicitOptionality: false })
  markAsSent(@path uid: uuid): NotificationUpdateResponse | ForbiddenError | NotFoundError;

  @doc("Mark notification as failed")
  @route("{uid}/failed")
  @patch(#{ implicitOptionality: false })
  markAsFailed(@path uid: uuid): NotificationUpdateResponse | ForbiddenError | NotFoundError;
}
