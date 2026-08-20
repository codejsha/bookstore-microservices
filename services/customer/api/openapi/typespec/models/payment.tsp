namespace CustomerService;

// ─── Payment models ────────────────────────────────────────────────────────

model PaymentFindResponse {
  uid: uuid;
  order_uid?: uuid;
  user_uid?: uuid;
  category_type?: string;
  currency?: string;
  amount?: float64;
  status?: string;
  reason?: string;
  @format("date-time") payment_date?: string;
}

model PaymentFindAllResponse is PaginatedResponse<PaymentFindResponse>;
