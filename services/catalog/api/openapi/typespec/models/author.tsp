namespace CatalogService;

// ─── Author models ──────────────────────────────────────────────────────────

model AuthorCreateRequest {
  name: string;
  bio?: string;
  birth_date?: string;
  death_date?: string;
  photo_uids?: string[];
  alternate_names?: string[];
  ol_key?: string;
}

model AuthorUpdateRequest {
  name?: string;
  bio?: string;
  birth_date?: string;
  death_date?: string;
  photo_uids?: string[];
  alternate_names?: string[];
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
