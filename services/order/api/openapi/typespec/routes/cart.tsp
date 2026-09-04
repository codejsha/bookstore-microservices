import "@typespec/http";

using Http;

namespace OrderService;

@route("/api/v1/cart")
@tag("Cart")
@useAuth(BearerAuth)
interface Cart {
  @doc("Get my cart") @get
  get(): CartFindResponse;

  @doc("Add item to cart") @route("items") @post
  addItem(@body body: CartAddItemRequest): CreatedResponse | BadRequestError | ConflictError;

  @doc("Update item quantity") @route("items/{uid}") @put
  updateItem(@path uid: uuid, @body body: CartUpdateItemRequest): CartItemFindResponse | BadRequestError | NotFoundError | ConflictError;

  @doc("Remove item from cart") @route("items/{uid}") @delete
  removeItem(@path uid: uuid): NoContentResponse | BadRequestError | NotFoundError;

  @doc("Clear cart") @delete
  clear(): NoContentResponse;

  @doc("Checkout cart to create order") @route("checkout") @post
  checkout(@body body: CartCheckoutRequest): OrderFindResponse | BadRequestError | ConflictError;
}
