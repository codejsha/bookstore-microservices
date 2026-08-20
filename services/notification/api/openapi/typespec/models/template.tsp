namespace NotificationService;

// ─── Template models ────────────────────────────────────────────────────────

model TemplateCreateRequest {
  notification_type: NotificationType;
  channel: Channel;
  title_template: string;
  content_template: string;
}

model TemplateUpdateRequest {
  notification_type?: NotificationType;
  channel?: Channel;
  title_template?: string;
  content_template?: string;
}

model TemplateItem {
  uid: uuid;
  notification_type: NotificationType;
  channel: Channel;
  title_template: string;
  content_template: string;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model TemplateFindResponse is TemplateItem;

model TemplateUpdateResponse is TemplateItem;

model TemplateFindAllResponse is PaginatedResponse<TemplateItem>;
