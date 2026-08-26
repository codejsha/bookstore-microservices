import "@typespec/http";

using Http;

namespace AdminService;

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
model UnauthorizedError {
  @statusCode statusCode: 401;
  @header contentType: "application/problem+json";
  type?: string;
  title: string;
  status: int32;
  detail?: string;
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

@error
model ConflictError {
  @statusCode statusCode: 409;
  @header contentType: "application/problem+json";
  type?: string;
  title: string;
  status: int32;
  detail?: string;
}

// ─── Scalars ────────────────────────────────────────────────────────────────

@doc("UUID string identifier")
scalar uuid extends string;

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
