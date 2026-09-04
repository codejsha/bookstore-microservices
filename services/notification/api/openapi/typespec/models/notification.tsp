namespace NotificationService;

// ─── Notification models ────────────────────────────────────────────────────

model SendNotificationRequest {
  user_uid: uuid;
  notification_type: NotificationType;
  channel: Channel;
  @minLength(1)
  @maxLength(255)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  title: string;
  @minLength(1)
  @maxLength(10000)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  content: string;
}

model NotificationItem {
  uid: uuid;
  user_uid: uuid;
  notification_type: NotificationType;
  channel: Channel;
  status: NotificationStatus;
  title: string;
  content: string;
  sent_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model NotificationFindResponse is NotificationItem;

model NotificationUpdateResponse is NotificationItem;

model NotificationFindAllResponse is PaginatedResponse<NotificationItem>;
