namespace IdentityService;

// ─── Enums ──────────────────────────────────────────────────────────────────

enum AuthRole {
  UNKNOWN,
  USER,
  STAFF,
  MANAGE,
  SYSTEM,
}

enum UserStatus {
  ACTIVE,
  SUSPENDED,
  DEACTIVATED,
}

// ─── User models ────────────────────────────────────────────────────────────

model UserRegisterRequest {
  @format("email")
  @maxLength(100)
  email: string;

  @secret
  password: string;

  @maxLength(50)
  first_name: string;

  @maxLength(50)
  last_name: string;

  @maxLength(30)
  phone?: string;

  roles?: AuthRole[];
}

model UserFindResponse {
  uid: uuid;

  @format("email")
  @maxLength(100)
  email: string;

  @maxLength(50)
  first_name: string;

  @maxLength(50)
  last_name: string;

  @maxLength(30)
  phone?: string;

  @minItems(1)
  roles: AuthRole[];

  status: UserStatus;
  last_login_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model UserUpdateRequest {
  @maxLength(50)
  first_name?: string;

  @maxLength(50)
  last_name?: string;

  @maxLength(30)
  phone?: string;
}

model UserRolesRequest {
  @minItems(1)
  roles: AuthRole[];
}

model UserFindAllResponse is PaginatedResponse<UserFindResponse>;
