namespace SupportService;

// ─── Category models ────────────────────────────────────────────────────────

model CategoryCreateRequest {
  @minLength(1)
  @maxLength(100)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  name: string;
  @maxLength(500)
  description?: string;
  parent_uid?: uuid;
}

model CategoryUpdateRequest {
  @minLength(1)
  @maxLength(100)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  name?: string;
  @maxLength(500)
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
