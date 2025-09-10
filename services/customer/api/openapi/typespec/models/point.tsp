namespace CustomerService;

// ─── Point change type ──────────────────────────────────────────────────────

enum PointChangeType {
  EARN,
  SPEND,
  ADJUST,
  EXPIRE,
}

// ─── Point models ───────────────────────────────────────────────────────────

model PointBalanceResponse {
  user_uid: uuid;
  balance: int32;
}

@doc("Earn points for a customer")
model PointEarnRequest {
  @minValue(1)
  amount: int32;

  reason?: string;
}

@doc("Spend points for a customer")
model PointSpendRequest {
  @minValue(1)
  amount: int32;

  reason?: string;
}

// ─── Point history models ───────────────────────────────────────────────────

model PointHistoryFindResponse {
  uid: uuid;
  change_type: PointChangeType;
  amount: int32;
  reason?: string;
  created_at: utcDateTime;
}

model PointHistoryFindAllResponse is PaginatedResponse<PointHistoryFindResponse>;
