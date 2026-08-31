namespace SupportService;

// ─── FAQ models ─────────────────────────────────────────────────────────────

model FaqCreateRequest {
  @minLength(1)
  @maxLength(500)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  question: string;
  @minLength(1)
  @maxLength(16000)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  answer: string;
  category_uid?: uuid;
  published?: boolean;
}

model FaqUpdateRequest {
  @minLength(1)
  @maxLength(500)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  question?: string;
  @minLength(1)
  @maxLength(16000)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  answer?: string;
  category_uid?: uuid;
  published?: boolean;
}

model FaqItem {
  uid: uuid;
  category_uid?: uuid;
  question: string;
  answer: string;
  view_count: int64;
  published: boolean;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model FaqFindResponse is FaqItem;

model FaqUpdateResponse is FaqItem;

model FaqFindAllResponse is PaginatedResponse<FaqItem>;
