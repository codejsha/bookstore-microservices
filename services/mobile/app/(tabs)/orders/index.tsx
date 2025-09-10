import { Card } from "@bookstore/design/native/card";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "expo-router";
import { Package, ShoppingCart } from "lucide-react-native";
import { useState } from "react";
import { FlatList, Pressable, RefreshControl, Text, View } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { type Order, orderListQueryOptions } from "@/domains/order";
import { ErrorState } from "@/shared/components/ErrorState";
import { SkeletonList } from "@/shared/components/Skeleton";
import { useThemeColor } from "@/shared/theme/useThemeColor";

const STATUS_FILTERS = [
  { label: "All", value: undefined },
  { label: "Pending", value: "PENDING" },
  { label: "Paid", value: "PAID" },
  { label: "Shipped", value: "SHIPPED" },
  { label: "Delivered", value: "DELIVERED" },
  { label: "Cancelled", value: "CANCELLED" },
];

function statusColor(status: Order["status"]): string {
  switch (status) {
    case "PENDING":
      return "bg-secondary text-secondary-foreground";
    case "PAID":
      return "bg-primary text-primary-foreground";
    case "SHIPPED":
      return "bg-accent text-accent-foreground";
    case "DELIVERED":
      return "bg-success text-primary-foreground";
    case "CANCELLED":
      return "bg-destructive text-destructive-foreground";
    default:
      return "bg-muted text-muted-foreground";
  }
}

export default function OrdersScreen() {
  const router = useRouter();
  const insets = useSafeAreaInsets();
  const { color } = useThemeColor();
  const [status, setStatus] = useState<string | undefined>();

  const { data, isLoading, isError, refetch, isRefetching } = useQuery(
    orderListQueryOptions({ status, size: 20, page: 0 }),
  );
  const items = data?.items ?? [];
  const total = data?.total ?? 0;

  return (
    <View className="flex-1 bg-background">
      <View className="border-b border-border px-5 pb-3 pt-4">
        <Text className="text-2xl font-bold text-foreground">Orders</Text>
        <Text className="mt-1 text-xs text-muted-foreground">
          {total > 0 ? `${total} order${total !== 1 ? "s" : ""}` : ""}
        </Text>
        <FlatList
          horizontal
          showsHorizontalScrollIndicator={false}
          className="mt-3"
          data={STATUS_FILTERS}
          keyExtractor={(f) => f.label}
          contentContainerClassName="gap-2"
          renderItem={({ item: f }) => {
            const active = status === f.value;
            return (
              <Text
                onPress={() => setStatus(f.value)}
                className={
                  active
                    ? "rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground"
                    : "rounded-md bg-muted px-3 py-1.5 text-sm font-medium text-muted-foreground"
                }
              >
                {f.label}
              </Text>
            );
          }}
        />
      </View>

      {isError ? (
        <ErrorState
          message="Couldn't load your orders."
          onRetry={() => refetch()}
        />
      ) : isLoading ? (
        <SkeletonList count={6} />
      ) : (
        <FlatList
          data={items}
          keyExtractor={(item) => item.uid}
          contentContainerClassName="gap-3 p-5"
          contentContainerStyle={{ paddingBottom: insets.bottom + 20 }}
          refreshControl={
            <RefreshControl
              refreshing={isRefetching}
              onRefresh={() => refetch()}
            />
          }
          ListEmptyComponent={
            <View className="items-center py-16">
              <ShoppingCart color={color("muted-foreground")} size={40} />
              <Text className="mt-2 text-base font-medium text-foreground">
                No orders found
              </Text>
              <Text className="mt-1 text-sm text-muted-foreground">
                {status
                  ? "Try a different filter."
                  : "You haven't placed any orders yet."}
              </Text>
            </View>
          }
          renderItem={({ item }) => (
            <Pressable
              onPress={() =>
                router.push({
                  pathname: "/(tabs)/orders/[uid]",
                  params: { uid: item.uid },
                })
              }
            >
              <Card className="flex-row items-center gap-3">
                <View className="h-10 w-10 items-center justify-center rounded-md bg-muted">
                  <Package color={color("muted-foreground")} size={20} />
                </View>
                <View className="flex-1">
                  <Text className="font-medium text-foreground">
                    {item.orderNumber}
                  </Text>
                  <Text className="mt-0.5 text-xs text-muted-foreground">
                    {new Date(item.createdAt).toLocaleDateString()}
                    {" · "}
                    {item.items?.length ?? 0} item
                    {(item.items?.length ?? 0) !== 1 ? "s" : ""}
                  </Text>
                </View>
                <View className="items-end gap-1">
                  <Text className="text-sm font-semibold text-foreground">
                    {item.currency} {item.totalAmount.toFixed(2)}
                  </Text>
                  <Text
                    className={`rounded px-2 py-0.5 text-[10px] font-semibold ${statusColor(item.status)}`}
                  >
                    {item.status}
                  </Text>
                </View>
              </Card>
            </Pressable>
          )}
        />
      )}
    </View>
  );
}
