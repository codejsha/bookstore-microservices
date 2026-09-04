namespace CatalogService;

// ─── Subject models ────────────────────────────────────────────────────────

model SubjectCreateRequest {
  @maxLength(255)
  name: string;
}

model SubjectUpdateRequest {
  @maxLength(255)
  name?: string;
}

model SubjectItem {
  uid: uuid;
  name: string;
}

model SubjectFindResponse is SubjectItem;

model SubjectUpdateResponse is SubjectItem;

model SubjectFindAllResponse is PaginatedResponse<SubjectItem>;
