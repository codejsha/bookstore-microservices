namespace NotificationService;

// ─── Notification models ────────────────────────────────────────────────────

model SendNotificationRequest {
  user_uid: uuid;
  notification_type: NotificationType;
  channel: Channel;
  title: string;
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
