import "@typespec/http";

using Http;

namespace AdminService;

// ─── Users ──────────────────────────────────────────────────────────────────

@doc("A user account as owned by identity. Roles are realm roles, kept in sync with Keycloak by identity.")
model AdminUserResponse {
  uid: uuid;

  @format("email")
  email: string;

  first_name: string;
  last_name: string;
  phone?: string;

  @doc("ACTIVE, SUSPENDED, or DEACTIVATED")
  status: string;

  @doc("Realm role names: USER, STAFF, MANAGE, SYSTEM")
  roles: string[];

  last_login_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model AdminUserFindAllResponse is PaginatedResponse<AdminUserResponse>;

@doc("Replaces the user's roles wholesale — a role left out of the list is revoked.")
model AdminUserRolesRequest {
  @doc("Realm role names. An unknown name is rejected rather than silently dropped.")
  roles: string[];
}
