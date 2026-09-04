namespace SupportService;

// ─── Ticket enums ───────────────────────────────────────────────────────────

enum TicketStatus {
  OPEN: "OPEN",
  IN_PROGRESS: "IN_PROGRESS",
  RESOLVED: "RESOLVED",
  CLOSED: "CLOSED",
}

enum TicketPriority {
  LOW: "LOW",
  MEDIUM: "MEDIUM",
  HIGH: "HIGH",
  URGENT: "URGENT",
}

// ─── Ticket models ──────────────────────────────────────────────────────────

model TicketCreateRequest {
  customer_uid: uuid;
  @minLength(1)
  @maxLength(255)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  subject: string;
  @minLength(1)
  @maxLength(16000)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  description: string;
  priority?: TicketPriority;
  category_uid?: uuid;
}

model TicketUpdateRequest {
  @minLength(1)
  @maxLength(255)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  subject?: string;
  @minLength(1)
  @maxLength(16000)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  description?: string;
  priority?: TicketPriority;
  category_uid?: uuid;
  assignee_uid?: uuid;
}

model TicketUpdateStatusRequest {
  status: TicketStatus;
}

model TicketItem {
  uid: uuid;
  customer_uid: uuid;
  category_uid?: uuid;
  assignee_uid?: uuid;
  subject: string;
  description: string;
  status: TicketStatus;
  priority: TicketPriority;
  resolved_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model TicketFindResponse is TicketItem;

model TicketUpdateResponse is TicketItem;

model TicketFindAllResponse is PaginatedResponse<TicketItem>;
