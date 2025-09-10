import "@typespec/http";

using Http;

namespace AdminService;

// ─── Reference data ─────────────────────────────────────────────────────────

@doc("An author, as referenced by a work")
model AdminAuthorItem {
  uid: uuid;
  name: string;
}

@doc("A subject (category) a work is filed under")
model AdminSubjectItem {
  uid: uuid;
  name: string;
}

model AdminAuthorFindAllResponse is PaginatedResponse<AdminAuthorItem>;

model AdminSubjectFindAllResponse is PaginatedResponse<AdminSubjectItem>;

// ─── Works ──────────────────────────────────────────────────────────────────

model AdminWorkResponse {
  uid: uuid;
  title: string;
  description?: string;

  @doc("Free-form year or date as recorded by the catalog, not a validated date")
  first_publish_date?: string;

  @doc("Open Library key, when the work was imported from there")
  ol_key?: string;

  authors: AdminAuthorItem[];
  subjects: AdminSubjectItem[];

  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model AdminWorkFindAllResponse is PaginatedResponse<AdminWorkResponse>;

model AdminWorkCreateRequest {
  title: string;
  description?: string;
  first_publish_date?: string;
  ol_key?: string;

  @doc("Authors are referenced by uid; they must already exist")
  author_uids: uuid[];

  @doc("Subjects are created on demand from their names, unlike authors")
  subject_names?: string[];
}

@doc("Every field is optional: an absent field leaves the stored value alone")
model AdminWorkUpdateRequest {
  title?: string;
  description?: string;
  first_publish_date?: string;
  ol_key?: string;
  author_uids?: uuid[];
  subject_names?: string[];
}
