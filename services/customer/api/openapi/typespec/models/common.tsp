import "@typespec/http";

using Http;

namespace CustomerService;

// ─── Error models ───────────────────────────────────────────────────────────

@error
model BadRequestError {
  @statusCode statusCode: 400;
  @header contentType: "application/problem+json";
  type?: string;
  title: string;
  status: int32;
  detail?: string;
  errors?: string[];
}

@error
model ForbiddenError {
  @statusCode statusCode: 403;
  @header contentType: "application/problem+json";
  type?: string;
  title: string;
  status: int32;
  detail?: string;
}

@error
model NotFoundError {
  @statusCode statusCode: 404;
  @header contentType: "application/problem+json";
  type?: string;
  title: string;
  status: int32;
  detail?: string;
}

// ─── Scalars ────────────────────────────────────────────────────────────────

@doc("UUID string identifier")
scalar uuid extends string;

// ─── Enums ──────────────────────────────────────────────────────────────────

enum AuthRole {
  UNKNOWN,
  SYSTEM,
  MANAGE,
  PROFILE,
  ORDER,
  VIEW,
}

enum OrderStatus {
  UNKNOWN,
  PENDING,
  PAID,
  SHIPPING,
  COMPLETED,
  CANCELLED,
}

// ─── Common query parameters ────────────────────────────────────────────────

alias Size = int32;
alias Page = int32;
alias Sort = string;

// ─── Common response models ─────────────────────────────────────────────────

model PaginatedResponse<T> {
  @minValue(0)
  total: int64;

  @minItems(0)
  items: T[];
}

model CreatedResponse {
  @statusCode statusCode: 201;
  @doc("Location of the created resource")
  @header Location: string;
}

model NoContentResponse {
  @statusCode statusCode: 204;
}
