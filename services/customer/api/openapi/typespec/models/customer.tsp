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

model CustomerUpdateRequest {
  @format("email") email?: string;
  @format("password") password?: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  roles?: AuthRole[];
}

model CustomerUpdateResponse is CustomerFindResponse;

model CustomerFindAllResponse is PaginatedResponse<CustomerFindResponse>;
