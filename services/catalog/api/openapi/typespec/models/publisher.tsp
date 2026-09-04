namespace CatalogService;

// ─── Publisher models ───────────────────────────────────────────────────────

model PublisherCreateRequest {
  @maxLength(255)
  name: string;

  @maxLength(255)
  address?: string;

  @maxLength(64)
  ol_key?: string;
}

model PublisherUpdateRequest {
  @maxLength(255)
  name?: string;

  @maxLength(255)
  address?: string;

  @maxLength(64)
  ol_key?: string;
}

model PublisherItem {
  uid: uuid;
  name: string;
  address?: string;
  ol_key?: string;
}

model PublisherFindResponse is PublisherItem;

model PublisherUpdateResponse is PublisherItem;

model PublisherFindAllResponse is PaginatedResponse<PublisherItem>;
