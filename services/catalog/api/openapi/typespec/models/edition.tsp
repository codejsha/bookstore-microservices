namespace CatalogService;

// ─── Edition models ────────────────────────────────────────────────────────

model EditionCreateRequest {
  @maxLength(512)
  title: string;

  @maxLength(10)
  isbn10?: string;

  @maxLength(13)
  isbn13?: string;

  number_of_pages?: int32;

  @maxLength(32)
  publish_date?: string;

  cover_uids?: string[];
  languages?: string[];

  @maxLength(64)
  physical_format?: string;

  description?: string;
  work_uid: uuid;
  publisher_uid?: uuid;

  @maxLength(64)
  ol_key?: string;
}

model EditionUpdateRequest {
  @maxLength(512)
  title?: string;

  @maxLength(10)
  isbn10?: string;

  @maxLength(13)
  isbn13?: string;

  number_of_pages?: int32;

  @maxLength(32)
  publish_date?: string;

  cover_uids?: string[];
  languages?: string[];

  @maxLength(64)
  physical_format?: string;

  description?: string;
  work_uid?: uuid;
  publisher_uid?: uuid;

  @maxLength(64)
  ol_key?: string;
}

model EditionItem {
  uid: uuid;
  title: string;
  isbn10?: string;
  isbn13?: string;
  number_of_pages?: int32;
  publish_date?: string;
  cover_uids?: string[];
  languages?: string[];
  physical_format?: string;
  description?: string;
  work: WorkItem;
  publisher?: PublisherItem;
  ol_key?: string;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model EditionFindResponse is EditionItem;

model EditionUpdateResponse is EditionItem;

model EditionFindAllResponse is PaginatedResponse<EditionItem>;
