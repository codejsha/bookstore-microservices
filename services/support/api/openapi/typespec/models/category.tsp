namespace SupportService;

// ─── Category models ────────────────────────────────────────────────────────

model CategoryCreateRequest {
  name: string;
  description?: string;
  parent_uid?: uuid;
}

model CategoryUpdateRequest {
  name?: string;
  description?: string;
  parent_uid?: uuid;
}

model CategoryItem {
  uid: uuid;
  name: string;
  description?: string;
  parent_uid?: uuid;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model CategoryFindResponse is CategoryItem;

model CategoryUpdateResponse is CategoryItem;

model CategoryFindAllResponse {
  @minItems(0)
  items: CategoryItem[];
}
