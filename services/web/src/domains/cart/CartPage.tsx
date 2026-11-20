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
import { useTranslation } from "react-i18next";
import { ErrorState } from "@/shared/components/ErrorState";
import { useToast } from "@/shared/components/toast";
import { useFormat } from "@/shared/hooks/useFormat";
import { useLocale } from "@/shared/i18n/useLocale";
import {
  useCheckout,
  useClearCart,
  useRemoveCartItem,
  useUpdateCartItem,
} from "./cart-mutations";
import { cartQueryOptions } from "./cart-queries";

const CART_CURRENCY = "USD";

export function CartPage() {
  const { t } = useTranslation(["cart", "common"]);
  const format = useFormat();
  const locale = useLocale();
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
          title={t("error.title")}
          message={t("error.message")}
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
        <h1 className="text-3xl font-bold tracking-tight">{t("title")}</h1>
        {!isEmpty && (
          <span className="text-muted-foreground text-lg">
            {t("itemCount", { count: totalItems })}
          </span>
        )}
      </div>

      {isEmpty ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16 space-y-4">
            <ShoppingCart className="size-16 text-muted-foreground" />
            <p className="text-xl font-medium text-muted-foreground">
              {t("empty")}
            </p>
            <Button render={<Link to="/$locale/books" params={{ locale }} />}>
              {t("browseBooks")}
            </Button>
          </CardContent>
        </Card>
      ) : (
        <>
          <ul aria-label={t("listLabel")} className="space-y-4">
            {items.map((item) => {
              const name =
                item.product_name ??
                t("common:product.fallbackName", { id: item.product_id });
              return (
                <li key={item.uid}>
                  <Card>
                    <CardContent className="flex items-center justify-between gap-4 py-4">
                      <div className="flex h-16 w-12 shrink-0 items-center justify-center rounded-md bg-muted">
                        <BookOpen className="size-5 text-muted-foreground" />
                      </div>

                      <div className="flex-1 min-w-0">
                        <p className="font-medium truncate">{name}</p>
                        <p className="text-sm text-muted-foreground">
                          {t("item.priceEach", {
                            price: format.price(item.price, CART_CURRENCY),
                          })}
                        </p>
                      </div>

                      <div className="flex items-center gap-2">
                        <Button
                          variant="outline"
                          size="icon"
                          className="size-8"
                          aria-label={t("item.decreaseQuantity", { name })}
                          disabled={updateItem.isPending}
                          onClick={() =>
                            item.quantity <= 1
                              ? removeItem.mutate(item.uid, {
                                  onSuccess: () =>
                                    toast.success(t("toast.itemRemoved")),
                                  onError: () =>
                                    toast.error(t("toast.removeFailed")),
                                })
                              : updateItem.mutate(
                                  {
                                    uid: item.uid,
                                    data: { quantity: item.quantity - 1 },
                                  },
                                  {
                                    onError: () =>
                                      toast.error(t("toast.updateFailed")),
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
                          aria-label={t("item.increaseQuantity", { name })}
                          disabled={updateItem.isPending}
                          onClick={() =>
                            updateItem.mutate(
                              {
                                uid: item.uid,
                                data: { quantity: item.quantity + 1 },
                              },
                              {
                                onError: () =>
                                  toast.error(t("toast.updateFailed")),
                              },
                            )
                          }
                        >
                          <Plus className="size-3.5" />
                        </Button>
                      </div>

                      <div className="w-24 text-right font-medium">
                        {format.price(item.subtotal, CART_CURRENCY)}
                      </div>

                      <Button
                        variant="ghost"
                        size="icon"
                        className="size-8 text-destructive hover:text-destructive"
                        aria-label={t("item.remove", { name })}
                        disabled={removeItem.isPending}
                        onClick={() =>
                          removeItem.mutate(item.uid, {
                            onSuccess: () =>
                              toast.success(t("toast.itemRemoved")),
                            onError: () => toast.error(t("toast.removeFailed")),
                          })
                        }
                      >
                        <Trash2 className="size-4" />
                      </Button>
                    </CardContent>
                  </Card>
                </li>
              );
            })}
          </ul>

          <Separator />
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">
                {t("summary.totalItems", { total: totalItems })}
              </p>
              <p className="text-2xl font-bold">
                {format.price(totalAmount, CART_CURRENCY)}
              </p>
            </div>
          </div>

          <div className="flex items-center justify-between">
            <Button
              variant="outline"
              disabled={clearCart.isPending}
              onClick={() =>
                clearCart.mutate(undefined, {
                  onSuccess: () => toast.success(t("toast.cleared")),
                  onError: () => toast.error(t("toast.clearFailed")),
                })
              }
            >
              {t("actions.clear")}
            </Button>
            <Button
              disabled={checkout.isPending}
              onClick={() =>
                checkout.mutate(
                  {
                    currency: CART_CURRENCY,
                    idempotency_key: crypto.randomUUID(),
                  },
                  {
                    onSuccess: () => {
                      toast.success(t("toast.orderPlaced"));
                      navigate({ to: "/$locale/orders", params: { locale } });
                    },
                    onError: () => toast.error(t("toast.checkoutFailed")),
                  },
                )
              }
            >
              {checkout.isPending
                ? t("actions.processing")
                : t("actions.checkout")}
            </Button>
          </div>
        </>
      )}
    </div>
  );
}
