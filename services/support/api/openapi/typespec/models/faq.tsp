namespace SupportService;

// ─── FAQ models ─────────────────────────────────────────────────────────────

model FaqCreateRequest {
  question: string;
  answer: string;
  category_uid?: uuid;
  published?: boolean;
}

model FaqUpdateRequest {
  question?: string;
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
