import { Button } from "@bookstore/design/native/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@bookstore/design/native/card";
import { useQuery } from "@tanstack/react-query";
import { useLocalSearchParams } from "expo-router";
import { MapPin, Package, Truck } from "lucide-react-native";
import { ScrollView, Text, View } from "react-native";
import { orderDetailQueryOptions, useCancelOrder } from "@/domains/order";
import { ErrorState } from "@/shared/components/ErrorState";
import { SkeletonCard } from "@/shared/components/Skeleton";
import { useThemeColor } from "@/shared/theme/useThemeColor";
import { useToast } from "@/shared/toast";

export default function OrderDetailScreen() {
  const { uid } = useLocalSearchParams<{ uid: string }>();
  const { color } = useThemeColor();
  const toast = useToast();
  const {
    data: order,
    isLoading,
    isError,
    refetch,
  } = useQuery(orderDetailQueryOptions(uid));
  const cancel = useCancelOrder();

  if (isError) {
    return (
      <ErrorState
        message="Couldn't load this order."
        onRetry={() => refetch()}
      />
    );
  }

  if (isLoading || !order) {
    return (
      <View className="flex-1 bg-background p-5">
        <SkeletonCard />
      </View>
    );
  }

  const canCancel = order.status === "PENDING";

  const handleCancel = () => {
    cancel.mutate(order.uid, {
      onSuccess: () => toast.success("Order cancelled"),
      onError: () => toast.error("Couldn't cancel order"),
    });
  };

  return (
    <ScrollView
      className="flex-1 bg-background"
      contentContainerClassName="gap-4 p-5"
    >
      <View>
        <Text className="text-2xl font-bold text-foreground">
          {order.orderNumber}
        </Text>
        <Text className="mt-1 text-sm text-muted-foreground">
          {new Date(order.createdAt).toLocaleString()}
          {"  ·  "}
          {order.status}
        </Text>
      </View>

      <Card>
        <CardHeader>
          <View className="flex-row items-center gap-2">
            <Package color={color("muted-foreground")} size={16} />
            <CardTitle>{`Items (${order.items?.length ?? 0})`}</CardTitle>
          </View>
        </CardHeader>
        <CardContent className="gap-3">
          {order.items?.map((item) => (
            <View
              key={item.uid}
              className="flex-row items-center justify-between"
            >
              <View className="flex-1">
                <Text className="text-sm font-medium text-foreground">
                  {item.productName ?? `Product #${item.productId}`}
                </Text>
                <Text className="mt-0.5 text-xs text-muted-foreground">
                  {item.quantity} × {item.currency} {item.price.toFixed(2)}
                </Text>
              </View>
              <Text className="text-sm font-semibold text-foreground">
                {item.currency} {item.subtotal.toFixed(2)}
              </Text>
            </View>
          ))}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Summary</CardTitle>
        </CardHeader>
        <CardContent className="gap-1.5">
          <SummaryRow
            label="Items"
            value={`${order.currency} ${order.itemsAmount.toFixed(2)}`}
          />
          {order.discountAmount > 0 ? (
            <SummaryRow
              label="Discount"
              value={`-${order.currency} ${order.discountAmount.toFixed(2)}`}
            />
          ) : null}
          <SummaryRow
            label="Shipping"
            value={`${order.currency} ${order.shippingAmount.toFixed(2)}`}
          />
          <SummaryRow
            label="Tax"
            value={`${order.currency} ${order.taxAmount.toFixed(2)}`}
          />
          <View className="mt-2 flex-row justify-between border-t border-border pt-2">
            <Text className="text-base font-semibold text-foreground">
              Total
            </Text>
            <Text className="text-base font-semibold text-foreground">
              {order.currency} {order.totalAmount.toFixed(2)}
            </Text>
          </View>
        </CardContent>
      </Card>

      {order.shipping ? (
        <Card>
          <CardHeader>
            <View className="flex-row items-center gap-2">
              <Truck color={color("muted-foreground")} size={16} />
              <CardTitle>Shipping</CardTitle>
            </View>
          </CardHeader>
          <CardContent>
            <View className="flex-row gap-2">
              <MapPin color={color("muted-foreground")} size={14} />
              <View className="flex-1">
                <Text className="text-sm font-medium text-foreground">
                  {order.shipping.recipientName}
                </Text>
                {order.shipping.recipientPhone ? (
                  <Text className="text-xs text-muted-foreground">
                    {order.shipping.recipientPhone}
                  </Text>
                ) : null}
                <Text className="text-xs text-muted-foreground">
                  {order.shipping.addressLine1}
                  {order.shipping.addressLine2
                    ? `, ${order.shipping.addressLine2}`
                    : ""}
                </Text>
                <Text className="text-xs text-muted-foreground">
                  {order.shipping.city}, {order.shipping.state}{" "}
                  {order.shipping.postalCode}, {order.shipping.country}
                </Text>
              </View>
            </View>
          </CardContent>
        </Card>
      ) : null}

      {canCancel ? (
        <Button
          variant="destructive"
          disabled={cancel.isPending}
          onPress={handleCancel}
        >
          {cancel.isPending ? "Cancelling..." : "Cancel order"}
        </Button>
      ) : null}
    </ScrollView>
  );
}

function SummaryRow({ label, value }: { label: string; value: string }) {
  return (
    <View className="flex-row justify-between">
      <Text className="text-sm text-muted-foreground">{label}</Text>
      <Text className="text-sm text-foreground">{value}</Text>
    </View>
  );
}
