namespace CatalogService;

// ─── Publisher models ───────────────────────────────────────────────────────

model PublisherCreateRequest {
  name: string;
  address?: string;
  ol_key?: string;
}

model PublisherUpdateRequest {
  name?: string;
  address?: string;
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
