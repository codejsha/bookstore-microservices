namespace CatalogService;

// ─── Subject models ────────────────────────────────────────────────────────

model SubjectCreateRequest {
  name: string;
}

model SubjectUpdateRequest {
  name?: string;
}

model SubjectItem {
  uid: uuid;
  name: string;
}

model SubjectFindResponse is SubjectItem;

model SubjectUpdateResponse is SubjectItem;

model SubjectFindAllResponse is PaginatedResponse<SubjectItem>;
