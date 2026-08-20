namespace IdentityService;

// ─── Enums ──────────────────────────────────────────────────────────────────

enum AuthRole {
  UNKNOWN,
  ADMIN,
  SYSTEM,
  MANAGE,
  PROFILE,
  ORDER,
  VIEW,
}

enum UserStatus {
  ACTIVE,
  SUSPENDED,
  DEACTIVATED,
}

// ─── User models ────────────────────────────────────────────────────────────

model UserRegisterRequest {
  @format("email")
  email: string;

  @secret
  password: string;

  first_name: string;
  last_name: string;
  phone?: string;
  roles?: AuthRole[];
}

model UserFindResponse {
  uid: uuid;

  @format("email")
  email: string;

  first_name: string;
  last_name: string;
  phone?: string;

  @minItems(1)
  roles: AuthRole[];

  status: UserStatus;
  last_login_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model UserUpdateRequest {
  first_name?: string;
  last_name?: string;
  phone?: string;
}

model UserRolesRequest {
  @minItems(1)
  roles: AuthRole[];
}

model UserFindAllResponse is PaginatedResponse<UserFindResponse>;
