namespace CustomerService;

// ─── Wishlist models ────────────────────────────────────────────────────────

model WishlistResponse {
  book_uids: uuid[];
}

@doc("Add books to wishlist")
model WishlistAddRequest {
  @minItems(1)
  book_uids: uuid[];
}

@doc("Remove books from wishlist")
model WishlistRemoveRequest {
  @minItems(1)
  book_uids: uuid[];
}
