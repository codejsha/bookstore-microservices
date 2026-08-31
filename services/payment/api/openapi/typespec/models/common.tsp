import "@typespec/http";

using Http;

namespace PaymentService;

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

@doc("ISO 4217 currency code (e.g. USD, KRW)")
@maxLength(3)
scalar CurrencyCode extends string;

// ─── Enums ──────────────────────────────────────────────────────────────────

enum PaymentStatus {
  requires_payment_method,
  requires_confirmation,
  requires_customer_action,
  requires_capture,
  processing,
  succeeded,
  failed,
  cancelled,
  partially_captured,
  expired,
}

enum CaptureMethod {
  automatic,
  manual,
  manual_multiple,
  scheduled,
}

enum AuthenticationType {
  three_ds,
  no_three_ds,
}

enum PaymentMethodEnum {
  card,
  wallet,
  bank_transfer,
  pay_later,
}

enum FutureUsage {
  on_session,
  off_session,
}

enum MandateType {
  single_use,
  multi_use,
}

enum MandateStatus {
  active,
  revoked,
  inactive,
  pending,
}

enum RefundStatus {
  @doc("Refund is pending")
  pending,

  @doc("Refund succeeded")
  succeeded,

  @doc("Refund failed")
  failed,

  @doc("Refund is under review")
  review,
}

enum RefundType {
  instant,
  scheduled,
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

// ─── Address model ──────────────────────────────────────────────────────────

model Address {
  city?: string;
  country?: string;
  line1?: string;
  line2?: string;
  state?: string;
  zip?: string;
  first_name?: string;
  last_name?: string;
}
