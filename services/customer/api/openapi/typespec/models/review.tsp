namespace CustomerService;

// ─── Review models ──────────────────────────────────────────────────────────

model ReviewCreateRequest {
  book_uid: uuid;

  @minValue(1)
  @maxValue(5)
  rating: int32;

  @maxLength(255)
  title?: string;

  content?: string;
}

model ReviewUpdateRequest {
  @minValue(1)
  @maxValue(5)
  rating?: int32;

  @maxLength(255)
  title?: string;

  content?: string;
}

model ReviewFindResponse {
  uid: uuid;
  user_uid: uuid;
  book_uid: uuid;

  @minValue(1)
  @maxValue(5)
  rating: int32;

  @maxLength(255)
  title?: string;

  content?: string;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model ReviewFindAllResponse is PaginatedResponse<ReviewFindResponse>;
