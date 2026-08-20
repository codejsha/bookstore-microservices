import { Badge } from "@bookstore/design/ui/badge";
import { Button } from "@bookstore/design/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@bookstore/design/ui/card";
import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { ChevronLeft, ChevronRight, Package, ShoppingCart } from "lucide-react";
import { useState } from "react";
import { ErrorState } from "@/shared/components/ErrorState";
import { orderListQueryOptions } from "./order-queries";
import type { Order } from "./types";
import { OrderStatus } from "./types";

const PAGE_SIZE = 10;

const STATUS_FILTERS = [
  { label: "All", value: undefined },
  { label: "Pending", value: "PENDING" },
  { label: "Paid", value: "PAID" },
  { label: "Shipped", value: "SHIPPED" },
  { label: "Delivered", value: "DELIVERED" },
  { label: "Cancelled", value: "CANCELLED" },
] as const;

function statusColor(status: OrderStatus): string {
  switch (status) {
    case OrderStatus.OrderStatusPending:
      return "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300";
    case OrderStatus.OrderStatusPaid:
      return "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300";
    case OrderStatus.OrderStatusShipped:
      return "bg-indigo-100 text-indigo-800 dark:bg-indigo-900 dark:text-indigo-300";
    case OrderStatus.OrderStatusDelivered:
      return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300";
    case OrderStatus.OrderStatusCancelled:
      return "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300";
    case OrderStatus.OrderStatusRefunded:
      return "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300";
    default:
      return "bg-muted text-muted-foreground";
  }
}

function OrderCard({ order }: { order: Order }) {
  return (
    <Link to="/orders/$uid" params={{ uid: order.uid }}>
      <Card className="transition-colors hover:bg-muted/50">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
                <Package className="size-5 text-muted-foreground" />
              </div>
              <div>
                <CardTitle className="text-sm font-medium">
                  {order.order_number}
                </CardTitle>
                <p className="text-xs text-muted-foreground">
                  {new Date(order.created_at).toLocaleDateString()}
                </p>
              </div>
            </div>
            <Badge className={statusColor(order.status)}>{order.status}</Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between text-sm">
            <span className="text-muted-foreground">
              {order.items?.length ?? 0} item
              {(order.items?.length ?? 0) !== 1 ? "s" : ""}
            </span>
            <span className="font-semibold">
              {order.currency} {order.total_amount.toFixed(2)}
            </span>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}

function Skeleton({ className }: { className?: string }) {
  return (
    <div className={`animate-pulse rounded-md bg-muted ${className ?? ""}`} />
  );
}

function OrderCardSkeleton() {
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Skeleton className="h-10 w-10 rounded-lg" />
            <div className="space-y-1.5">
              <Skeleton className="h-4 w-28" />
              <Skeleton className="h-3 w-20" />
            </div>
          </div>
          <Skeleton className="h-5 w-16 rounded-full" />
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between">
          <Skeleton className="h-4 w-16" />
          <Skeleton className="h-4 w-20" />
        </div>
      </CardContent>
    </Card>
  );
}

export function OrderList() {
  const [status, setStatus] = useState<string | undefined>();
  const [page, setPage] = useState(0);

  const {
    data: orders,
    isLoading,
    isFetching,
    isError,
    refetch,
  } = useQuery(orderListQueryOptions({ status, size: PAGE_SIZE, page }));

  const total = orders?.total ?? 0;
  const items = orders?.items ?? [];
  const totalPages = Math.ceil(total / PAGE_SIZE);

  return (
    <div className="container mx-auto space-y-6 px-4 py-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Orders</h1>
          <p className="text-sm text-muted-foreground min-h-5">
            {total > 0 && `${total} order${total !== 1 ? "s" : ""}`}
          </p>
        </div>
      </div>

      <div className="flex flex-wrap gap-2">
        {STATUS_FILTERS.map((f) => (
          <button
            key={f.label}
            type="button"
            onClick={() => {
              setStatus(f.value);
              setPage(0);
            }}
            className={`inline-flex items-center rounded-md px-3 py-1 text-sm font-medium transition-colors ${
              status === f.value
                ? "bg-primary text-primary-foreground"
                : "bg-muted text-muted-foreground hover:text-foreground"
            }`}
          >
            {f.label}
          </button>
        ))}
      </div>

      <div className="relative min-h-[300px]">
        <div
          className={`grid gap-5 transition-opacity duration-300 ${isLoading ? "opacity-100" : "opacity-0 absolute inset-0 pointer-events-none"}`}
        >
          {Array.from({ length: PAGE_SIZE }).map((_, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: skeleton placeholders
            <OrderCardSkeleton key={i} />
          ))}
        </div>

        {!isLoading && isError && (
          <ErrorState
            title="Couldn't load orders"
            message="There was a problem loading your orders. Please try again."
            onRetry={() => refetch()}
            showHome={false}
          />
        )}

        {!isLoading && !isError && (
          <div
            className={`transition-opacity duration-200 ${isFetching ? "opacity-40" : "opacity-100"}`}
          >
            {items.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-20 text-center animate-in fade-in duration-300">
                <ShoppingCart className="size-12 text-muted-foreground mb-4" />
                <p className="text-lg font-medium">No orders found</p>
                <p className="mt-1 text-sm text-muted-foreground">
                  {status
                    ? "Try a different filter."
                    : "You haven't placed any orders yet."}
                </p>
              </div>
            ) : (
              <div className="grid gap-5 animate-in fade-in duration-300">
                {items.map((order) => (
                  <OrderCard key={order.uid} order={order} />
                ))}
              </div>
            )}
          </div>
        )}
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="icon-sm"
            disabled={page === 0}
            onClick={() => setPage((p) => p - 1)}
          >
            <ChevronLeft />
          </Button>
          <span className="text-sm text-muted-foreground">
            Page {page + 1} of {totalPages}
          </span>
          <Button
            variant="outline"
            size="icon-sm"
            disabled={page >= totalPages - 1}
            onClick={() => setPage((p) => p + 1)}
          >
            <ChevronRight />
          </Button>
        </div>
      )}
    </div>
  );
}
