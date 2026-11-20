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
import { useTranslation } from "react-i18next";
import { ErrorState } from "@/shared/components/ErrorState";
import { useFormat } from "@/shared/hooks/useFormat";
import { useLocale } from "@/shared/i18n/useLocale";
import { orderListQueryOptions } from "./order-queries";
import { useOrderStatusLabel } from "./order-status";
import type { Order } from "./types";
import { OrderStatus } from "./types";

const PAGE_SIZE = 10;

const STATUS_FILTERS = [
  { labelKey: "list.filter.all", value: undefined },
  { labelKey: "list.filter.pending", value: "PENDING" },
  { labelKey: "list.filter.paid", value: "PAID" },
  { labelKey: "list.filter.shipped", value: "SHIPPED" },
  { labelKey: "list.filter.delivered", value: "DELIVERED" },
  { labelKey: "list.filter.cancelled", value: "CANCELLED" },
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
  const { t } = useTranslation("order");
  const format = useFormat();
  const locale = useLocale();
  const statusLabel = useOrderStatusLabel();

  return (
    <Link to="/$locale/orders/$uid" params={{ locale, uid: order.uid }}>
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
                  {format.date(order.created_at)}
                </p>
              </div>
            </div>
            <Badge className={statusColor(order.status)}>
              {statusLabel(order.status)}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between text-sm">
            <span className="text-muted-foreground">
              {t("list.itemCount", { count: order.items?.length ?? 0 })}
            </span>
            <span className="font-semibold">
              {format.money(order.total_amount, order.currency)}
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
  const { t } = useTranslation("order");
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
          <h1 className="text-2xl font-bold tracking-tight">
            {t("list.title")}
          </h1>
          <p className="text-sm text-muted-foreground min-h-5">
            {total > 0 && t("list.count", { count: total })}
          </p>
        </div>
      </div>

      <div className="flex flex-wrap gap-2">
        {STATUS_FILTERS.map((f) => (
          <button
            key={f.labelKey}
            type="button"
            aria-pressed={status === f.value}
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
            {t(f.labelKey)}
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
            title={t("list.errorTitle")}
            message={t("list.errorMessage")}
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
                <p className="text-lg font-medium">{t("list.emptyTitle")}</p>
                <p className="mt-1 text-sm text-muted-foreground">
                  {status ? t("list.emptyFiltered") : t("list.emptyAll")}
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
            aria-label={t("list.previousPage")}
            disabled={page === 0}
            onClick={() => setPage((p) => p - 1)}
          >
            <ChevronLeft />
          </Button>
          <span className="text-sm text-muted-foreground">
            {t("list.pageOf", { page: page + 1, totalPages })}
          </span>
          <Button
            variant="outline"
            size="icon-sm"
            aria-label={t("list.nextPage")}
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
