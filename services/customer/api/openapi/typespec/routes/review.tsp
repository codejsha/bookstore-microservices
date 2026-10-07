import "@typespec/http";

using Http;

namespace CustomerService;

// ─── My Reviews (customer perspective) ──────────────────────────────────────

@route("/api/v1/customers/{uid}/reviews")
@tag("Review")
@useAuth(BearerAuth)
interface CustomerReviews {
  @doc("Get all reviews by a customer")
  @get
  getAll(
    @path uid: uuid,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): ReviewFindAllResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Write a review for a book")
  @post
  write(
    @path uid: uuid,
    @body body: ReviewCreateRequest,
  ):
    | CreatedResponse
    | BadRequestError
    | NotFoundError
    | UnauthorizedError
    | ForbiddenError
    | ConflictError;

  @doc("Get a specific review")
  @route("{review_uid}")
  @get
  read(
    @path uid: uuid,
    @path review_uid: uuid,
  ): ReviewFindResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Edit a review")
  @route("{review_uid}")
  @put
  edit(
    @path uid: uuid,
    @path review_uid: uuid,
    @body body: ReviewUpdateRequest,
  ):
    | ReviewFindResponse
    | BadRequestError
    | NotFoundError
    | UnauthorizedError
    | ForbiddenError;

  @doc("Delete a review")
  @route("{review_uid}")
  @delete
  remove(
    @path uid: uuid,
    @path review_uid: uuid,
  ): NoContentResponse | NotFoundError | UnauthorizedError | ForbiddenError;
}

// ─── Book Reviews (public, read-only) ───────────────────────────────────────

@route("/api/v1/books/{book_uid}/reviews")
@tag("Review")
interface BookReviews {
  @doc("Get all reviews for a book")
  @get
  getAll(
    @path book_uid: uuid,
    @query rating?: int32,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): ReviewFindAllResponse | BadRequestError;
}
