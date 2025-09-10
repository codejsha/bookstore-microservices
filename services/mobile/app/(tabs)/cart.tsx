import { Button } from "@bookstore/design/native/button";
import { Card } from "@bookstore/design/native/card";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "expo-router";
import {
  BookOpen,
  Minus,
  Plus,
  ShoppingCart,
  Trash2,
} from "lucide-react-native";
import {
  Pressable,
  RefreshControl,
  ScrollView,
  Text,
  View,
} from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import {
  cartQueryOptions,
  useCheckout,
  useClearCart,
  useRemoveCartItem,
  useUpdateCartItem,
} from "@/domains/cart";
import { ErrorState } from "@/shared/components/ErrorState";
import { SkeletonList } from "@/shared/components/Skeleton";
import { useThemeColor } from "@/shared/theme/useThemeColor";
import { useToast } from "@/shared/toast";

export default function CartScreen() {
  const router = useRouter();
  const insets = useSafeAreaInsets();
  const { color } = useThemeColor();
  const toast = useToast();
  const {
    data: cart,
    isLoading,
    isError,
    refetch,
    isRefetching,
  } = useQuery(cartQueryOptions());
  const update = useUpdateCartItem();
  const remove = useRemoveCartItem();
  const clear = useClearCart();
  const checkout = useCheckout();

  if (isError) {
    return (
      <ErrorState
        message="Couldn't load your cart."
        onRetry={() => refetch()}
      />
    );
  }

  if (isLoading) {
    return (
      <View className="flex-1 bg-background">
        <SkeletonList count={4} />
      </View>
    );
  }

  const items = cart?.items ?? [];
  const isEmpty = items.length === 0;

  const handleUpdate = (uid: string, quantity: number) => {
    update.mutate(
      { uid, data: { quantity } },
      { onError: () => toast.error("Couldn't update quantity") },
    );
  };

  const handleRemove = (uid: string) => {
    remove.mutate(uid, {
      onSuccess: () => toast.success("Removed from cart"),
      onError: () => toast.error("Couldn't remove item"),
    });
  };

  if (isEmpty) {
    return (
      <View className="flex-1 items-center justify-center bg-background px-5">
        <ShoppingCart color={color("muted-foreground")} size={48} />
        <Text className="mt-3 text-base font-medium text-foreground">
          Your cart is empty
        </Text>
        <View className="mt-4">
          <Button onPress={() => router.push("/(tabs)/books")}>
            Browse books
          </Button>
        </View>
      </View>
    );
  }

  return (
    <View className="flex-1 bg-background">
      <ScrollView
        contentContainerClassName="gap-3 p-5"
        refreshControl={
          <RefreshControl
            refreshing={isRefetching}
            onRefresh={() => refetch()}
          />
        }
      >
        {items.map((item) => (
          <Card key={item.uid} className="flex-row items-center gap-3">
            <View className="h-16 w-12 items-center justify-center rounded-md bg-muted">
              <BookOpen color={color("muted-foreground")} size={20} />
            </View>
            <View className="flex-1">
              <Text className="font-medium text-foreground" numberOfLines={2}>
                {item.product_name ?? `Product #${item.product_id}`}
              </Text>
              <Text className="mt-1 text-xs text-muted-foreground">
                ${item.price.toFixed(2)} each
              </Text>
              <View className="mt-2 flex-row items-center gap-2">
                <Pressable
                  onPress={() =>
                    item.quantity <= 1
                      ? handleRemove(item.uid)
                      : handleUpdate(item.uid, item.quantity - 1)
                  }
                  className="h-7 w-7 items-center justify-center rounded-md border border-border"
                >
                  <Minus color={color("muted-foreground")} size={14} />
                </Pressable>
                <Text className="w-8 text-center text-sm font-medium text-foreground">
                  {item.quantity}
                </Text>
                <Pressable
                  onPress={() => handleUpdate(item.uid, item.quantity + 1)}
                  className="h-7 w-7 items-center justify-center rounded-md border border-border"
                >
                  <Plus color={color("muted-foreground")} size={14} />
                </Pressable>
              </View>
            </View>
            <View className="items-end">
              <Text className="text-sm font-semibold text-foreground">
                ${item.subtotal.toFixed(2)}
              </Text>
              <Pressable
                onPress={() => handleRemove(item.uid)}
                className="mt-2 h-7 w-7 items-center justify-center"
              >
                <Trash2 color={color("destructive")} size={16} />
              </Pressable>
            </View>
          </Card>
        ))}
      </ScrollView>

      <View
        className="border-t border-border bg-background p-5"
        style={{ paddingBottom: insets.bottom + 20 }}
      >
        <View className="mb-3 flex-row items-center justify-between">
          <Text className="text-sm text-muted-foreground">
            Total ({cart?.total_items} items)
          </Text>
          <Text className="text-xl font-bold text-foreground">
            ${cart?.total_amount.toFixed(2)}
          </Text>
        </View>
        <View className="flex-row gap-2">
          <View className="flex-1">
            <Button
              variant="outline"
              disabled={clear.isPending}
              onPress={() =>
                clear.mutate(undefined, {
                  onSuccess: () => toast.success("Cart cleared"),
                  onError: () => toast.error("Couldn't clear cart"),
                })
              }
            >
              Clear
            </Button>
          </View>
          <View className="flex-1">
            <Button
              disabled={checkout.isPending}
              onPress={() =>
                checkout.mutate(
                  {
                    currency: "USD",
                    idempotency_key: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
                  },
                  {
                    onSuccess: () => {
                      toast.success("Order placed");
                      router.push("/(tabs)/orders");
                    },
                    onError: () => toast.error("Checkout failed"),
                  },
                )
              }
            >
              {checkout.isPending ? "..." : "Checkout"}
            </Button>
          </View>
        </View>
      </View>
    </View>
  );
}
