import { Button } from "@bookstore/design/ui/button";
import { Card, CardContent } from "@bookstore/design/ui/card";
import { Separator } from "@bookstore/design/ui/separator";
import { useQuery } from "@tanstack/react-query";
import { Link, useNavigate } from "@tanstack/react-router";
import {
  BookOpen,
  Loader2,
  Minus,
  Plus,
  ShoppingCart,
  Trash2,
} from "lucide-react";
import { ErrorState } from "@/shared/components/ErrorState";
import { useToast } from "@/shared/components/toast";
import {
  useCheckout,
  useClearCart,
  useRemoveCartItem,
  useUpdateCartItem,
} from "./cart-mutations";
import { cartQueryOptions } from "./cart-queries";

export function CartPage() {
  const {
    data: cart,
    isLoading,
    isError,
    refetch,
  } = useQuery(cartQueryOptions());
  const navigate = useNavigate();
  const toast = useToast();

  const updateItem = useUpdateCartItem();
  const removeItem = useRemoveCartItem();
  const clearCart = useClearCart();
  const checkout = useCheckout();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="size-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="container mx-auto max-w-4xl px-4">
        <ErrorState
          title="Couldn't load your cart"
          message="There was a problem loading your cart. Please try again."
          onRetry={() => refetch()}
        />
      </div>
    );
  }

  const items = cart?.items ?? [];
  const totalItems = cart?.total_items ?? 0;
  const totalAmount = cart?.total_amount ?? 0;
  const isEmpty = items.length === 0;

  return (
    <div className="container mx-auto max-w-4xl px-4 py-8 space-y-8">
      <div className="flex items-center gap-3">
        <ShoppingCart className="size-7" />
        <h1 className="text-3xl font-bold tracking-tight">Shopping Cart</h1>
        {!isEmpty && (
          <span className="text-muted-foreground text-lg">
            ({totalItems} {totalItems === 1 ? "item" : "items"})
          </span>
        )}
      </div>

      {isEmpty ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16 space-y-4">
            <ShoppingCart className="size-16 text-muted-foreground" />
            <p className="text-xl font-medium text-muted-foreground">
              Your cart is empty
            </p>
            <Button render={<Link to="/books" />}>Browse Books</Button>
          </CardContent>
        </Card>
      ) : (
        <>
          <div className="space-y-4">
            {items.map((item) => (
              <Card key={item.uid}>
                <CardContent className="flex items-center justify-between gap-4 py-4">
                  <div className="flex h-16 w-12 shrink-0 items-center justify-center rounded-md bg-muted">
                    <BookOpen className="size-5 text-muted-foreground" />
                  </div>

                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate">
                      {item.product_name ?? `Product #${item.product_id}`}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      ${item.price.toFixed(2)} each
                    </p>
                  </div>

                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="icon"
                      className="size-8"
                      disabled={updateItem.isPending}
                      onClick={() =>
                        item.quantity <= 1
                          ? removeItem.mutate(item.uid, {
                              onSuccess: () => toast.success("Item removed"),
                              onError: () =>
                                toast.error("Couldn't remove item"),
                            })
                          : updateItem.mutate(
                              {
                                uid: item.uid,
                                data: { quantity: item.quantity - 1 },
                              },
                              {
                                onError: () =>
                                  toast.error("Couldn't update quantity"),
                              },
                            )
                      }
                    >
                      <Minus className="size-3.5" />
                    </Button>
                    <span className="w-8 text-center text-sm font-medium">
                      {item.quantity}
                    </span>
                    <Button
                      variant="outline"
                      size="icon"
                      className="size-8"
                      disabled={updateItem.isPending}
                      onClick={() =>
                        updateItem.mutate(
                          {
                            uid: item.uid,
                            data: { quantity: item.quantity + 1 },
                          },
                          {
                            onError: () =>
                              toast.error("Couldn't update quantity"),
                          },
                        )
                      }
                    >
                      <Plus className="size-3.5" />
                    </Button>
                  </div>

                  <div className="w-24 text-right font-medium">
                    ${item.subtotal.toFixed(2)}
                  </div>

                  <Button
                    variant="ghost"
                    size="icon"
                    className="size-8 text-destructive hover:text-destructive"
                    disabled={removeItem.isPending}
                    onClick={() =>
                      removeItem.mutate(item.uid, {
                        onSuccess: () => toast.success("Item removed"),
                        onError: () => toast.error("Couldn't remove item"),
                      })
                    }
                  >
                    <Trash2 className="size-4" />
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>

          <Separator />
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">
                Total items: {totalItems}
              </p>
              <p className="text-2xl font-bold">${totalAmount.toFixed(2)}</p>
            </div>
          </div>

          <div className="flex items-center justify-between">
            <Button
              variant="outline"
              disabled={clearCart.isPending}
              onClick={() =>
                clearCart.mutate(undefined, {
                  onSuccess: () => toast.success("Cart cleared"),
                  onError: () => toast.error("Couldn't clear cart"),
                })
              }
            >
              Clear Cart
            </Button>
            <Button
              disabled={checkout.isPending}
              onClick={() =>
                checkout.mutate(
                  { currency: "USD", idempotency_key: crypto.randomUUID() },
                  {
                    onSuccess: () => {
                      toast.success("Order placed");
                      navigate({ to: "/orders" });
                    },
                    onError: () => toast.error("Checkout failed"),
                  },
                )
              }
            >
              {checkout.isPending ? "Processing..." : "Checkout"}
            </Button>
          </div>
        </>
      )}
    </div>
  );
}
