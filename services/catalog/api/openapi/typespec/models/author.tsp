namespace CatalogService;

// ─── Author models ──────────────────────────────────────────────────────────

model AuthorCreateRequest {
  @maxLength(255)
  name: string;

  bio?: string;

  @maxLength(32)
  birth_date?: string;

  @maxLength(32)
  death_date?: string;

  photo_uids?: string[];
  alternate_names?: string[];

  @maxLength(64)
  ol_key?: string;
}

model AuthorUpdateRequest {
  @maxLength(255)
  name?: string;

  bio?: string;

  @maxLength(32)
  birth_date?: string;

  @maxLength(32)
  death_date?: string;

  photo_uids?: string[];
  alternate_names?: string[];

  @maxLength(64)
  ol_key?: string;
}

model AuthorItem {
  uid: uuid;
  name: string;
  bio?: string;
  birth_date?: string;
  death_date?: string;
  photo_uids?: string[];
  alternate_names?: string[];
  ol_key?: string;
}

model AuthorFindResponse is AuthorItem;

model AuthorUpdateResponse is AuthorItem;

model AuthorFindAllResponse is PaginatedResponse<AuthorItem>;
