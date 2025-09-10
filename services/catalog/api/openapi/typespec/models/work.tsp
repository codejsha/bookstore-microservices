namespace CatalogService;

// ─── Work models ───────────────────────────────────────────────────────────

model WorkCreateRequest {
  title: string;
  description?: string;
  cover_uids?: string[];
  first_publish_date?: string;
  ol_key?: string;

  @minItems(1)
  author_uids: uuid[];

  subject_names?: string[];
}

model WorkUpdateRequest {
  title?: string;
  description?: string;
  cover_uids?: string[];
  first_publish_date?: string;
  ol_key?: string;
  author_uids?: uuid[];
  subject_names?: string[];
}

model WorkItem {
  uid: uuid;
  title: string;
  description?: string;
  cover_uids?: string[];
  first_publish_date?: string;
  ol_key?: string;
  authors: AuthorItem[];
  subjects: SubjectItem[];
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model WorkFindResponse is WorkItem;

model WorkUpdateResponse is WorkItem;

model WorkFindAllResponse is PaginatedResponse<WorkItem>;

// ─── Full-text search ──────────────────────────────────────────────────────

model WorkSearchHighlight {
  title?: string[];
  description?: string[];
}

model WorkSearchItem {
  uid: uuid;
  title: string;
  description?: string;
  cover_uids?: string[];
  first_publish_date?: string;
  ol_key?: string;
  authors: AuthorItem[];
  subjects: SubjectItem[];
  score: float32;
  highlight?: WorkSearchHighlight;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model WorkSearchResponse is PaginatedResponse<WorkSearchItem>;
