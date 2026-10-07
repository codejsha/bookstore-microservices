import "@typespec/http";

using Http;

namespace CustomerService;

@route("/api/v1/customers/{uid}/wishlist")
@tag("Wishlist")
@useAuth(BearerAuth)
interface Wishlists {
  @doc("Get customer's wishlist")
  @get
  get(@path uid: uuid): WishlistResponse | UnauthorizedError | ForbiddenError;

  @doc("Add books to customer's wishlist")
  @route("add")
  @post
  addBooks(@path uid: uuid, @body body: WishlistAddRequest):
    | WishlistResponse
    | BadRequestError
    | UnauthorizedError
    | ForbiddenError;

  @doc("Remove books from customer's wishlist")
  @route("remove")
  @post
  removeBooks(@path uid: uuid, @body body: WishlistRemoveRequest):
    | WishlistResponse
    | BadRequestError
    | UnauthorizedError
    | ForbiddenError;
}
