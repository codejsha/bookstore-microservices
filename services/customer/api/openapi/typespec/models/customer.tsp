namespace CustomerService;

// ─── User model ────────────────────────────────────────────────────────────

model User {
  uid: uuid;
  @format("email") email?: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
}

// ─── Customer models ───────────────────────────────────────────────────────

model CustomerFindResponse {
  uid: uuid;
  @format("email") email?: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  roles?: AuthRole[];
}

model CustomerFindAllResponse is PaginatedResponse<CustomerFindResponse>;
