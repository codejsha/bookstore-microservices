import "@typespec/http";

using Http;

namespace CatalogService;

// ─── Error models ───────────────────────────────────────────────────────────

@error
model BadRequestError {
  @statusCode statusCode: 400;
  code: int32;
  message: string;
  details?: string[];
}

@error
model ForbiddenError {
  @statusCode statusCode: 403;
  code: int32;
  message: string;
}

@error
model NotFoundError {
  @statusCode statusCode: 404;
  code: int32;
  message: string;
}

// ─── Scalars ────────────────────────────────────────────────────────────────

@doc("UUid string identifier")
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

model CreatedResponse {
  @statusCode statusCode: 201;
  @doc("Location of the created resource")
  @header Location: string;
}

model NoContentResponse {
  @statusCode statusCode: 204;
}
