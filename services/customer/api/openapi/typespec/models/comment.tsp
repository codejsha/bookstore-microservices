namespace CustomerService;

// ─── Comment enums ──────────────────────────────────────────────────────────

enum AuthorRole {
  CUSTOMER,
  AGENT,
}

// ─── Comment models ─────────────────────────────────────────────────────────

model CommentCreateRequest {
  @minLength(1)
  @maxLength(16000)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  body: string;

  internal?: boolean;
}

model CommentFindResponse {
  uid: uuid;
  ticket_uid: uuid;
  author_uid: uuid;
  author_role: AuthorRole;
  body: string;
  internal: boolean;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model CommentFindAllResponse {
  @minItems(0)
  items: CommentFindResponse[];
}
