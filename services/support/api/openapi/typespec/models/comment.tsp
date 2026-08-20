namespace SupportService;

// ─── Comment enums ──────────────────────────────────────────────────────────

enum AuthorRole {
  CUSTOMER: "CUSTOMER",
  AGENT: "AGENT",
}

// ─── Comment models ─────────────────────────────────────────────────────────

// Authorship (author_uid / author_role) is derived server-side from the
// authenticated principal — never client-supplied — so it is intentionally
// absent from the create request.
model CommentCreateRequest {
  body: string;
  internal?: boolean;
}

model CommentItem {
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
  items: CommentItem[];
}
