import "@typespec/protobuf";

using TypeSpec.Protobuf;

@package({
  name: "payment.v1",
})
namespace ProtoPayment;

// ─── Enum ───────────────────────────────────────────────────────────────────

enum PaymentStatus {
  PAYMENT_STATUS_UNSPECIFIED: 0,
  PAYMENT_STATUS_REQUIRES_PAYMENT_METHOD: 1,
  PAYMENT_STATUS_REQUIRES_CONFIRMATION: 2,
  PAYMENT_STATUS_REQUIRES_CUSTOMER_ACTION: 3,
  PAYMENT_STATUS_REQUIRES_CAPTURE: 4,
  PAYMENT_STATUS_PROCESSING: 5,
  PAYMENT_STATUS_SUCCEEDED: 6,
  PAYMENT_STATUS_FAILED: 7,
  PAYMENT_STATUS_CANCELLED: 8,
  PAYMENT_STATUS_EXPIRED: 9,
  PAYMENT_STATUS_PARTIALLY_CAPTURED: 10,
}

// ─── Entity ─────────────────────────────────────────────────────────────────

model Payment {
  @field(1) uid: string;
  @field(2) payment_uid: string;
  @field(3) customer_uid?: string;
  @field(4) amount: int64;
  @field(5) currency: string;
  @field(6) status: PaymentStatus;
  @field(7) payment_method?: string;
  @field(8) error_code?: string;
  @field(9) error_message?: string;
}

// ─── Messages ───────────────────────────────────────────────────────────────

model ListPaymentsRequest {
  @field(1) customer_uid?: string;
  @field(2) status?: PaymentStatus;
  @field(3) page_size?: int32;
  @field(4) page_token?: string;
  @field(5) order_by?: string;
}

model ListPaymentsResponse {
  @field(1) payments: Payment[];
  @field(2) next_page_token?: string;
  @field(3) total_size: int32;
}

model FindPaymentRequest {
  @field(1) uid: string;
}

model FindPaymentResponse {
  @field(1) payment: Payment;
}

// ─── Refund ─────────────────────────────────────────────────────────────────

enum RefundStatus {
  REFUND_STATUS_UNSPECIFIED: 0,
  REFUND_STATUS_PENDING: 1,
  REFUND_STATUS_SUCCEEDED: 2,
  REFUND_STATUS_FAILED: 3,
  REFUND_STATUS_REVIEW: 4,
}

model Refund {
  @field(1) uid: string;
  @field(2) payment_id: string;
  @field(3) amount: int64;
  @field(4) currency: string;
  @field(5) status: RefundStatus;
  @field(6) reason?: string;
  @field(7) refund_type: string;
  @field(8) error_code?: string;
  @field(9) error_message?: string;
  @field(10) created_at: string;
}

model ListRefundsRequest {
  @field(1) payment_id?: string;
  @field(2) status?: RefundStatus;
  @field(3) page_size?: int32;
  @field(4) page_token?: string;
  @field(5) order_by?: string;
}

model ListRefundsResponse {
  @field(1) refunds: Refund[];
  @field(2) next_page_token?: string;
  @field(3) total_size: int32;
}

model FindRefundRequest {
  @field(1) uid: string;
}

model FindRefundResponse {
  @field(1) refund: Refund;
}

// ─── Service ────────────────────────────────────────────────────────────────

@Protobuf.service
interface PaymentService {
  ListPayments(...ListPaymentsRequest): ListPaymentsResponse;
  FindPayment(...FindPaymentRequest): FindPaymentResponse;
  ListRefunds(...ListRefundsRequest): ListRefundsResponse;
  FindRefund(...FindRefundRequest): FindRefundResponse;
}
