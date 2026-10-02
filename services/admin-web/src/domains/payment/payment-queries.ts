import { queryOptions } from "@tanstack/react-query";
import { api, retryUnlessNotFound } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import { toQueryString } from "@/shared/api/query-string";
import type { Paged } from "@/shared/api/types";
import type {
  Payment,
  PaymentListParams,
  Refund,
  RefundListParams,
} from "./types";

export const paymentKeys = {
  all: ["payment"] as const,
  list: (params: PaymentListParams) =>
    [...paymentKeys.all, "list", params] as const,
  detail: (uid: string) => [...paymentKeys.all, "detail", uid] as const,
  refundList: (params: RefundListParams) =>
    [...paymentKeys.all, "refund-list", params] as const,
  refund: (uid: string) => [...paymentKeys.all, "refund", uid] as const,
};

export const paymentsQueryOptions = (params: PaymentListParams) =>
  queryOptions({
    queryKey: paymentKeys.list(params),
    queryFn: () =>
      api.get<Paged<Payment>>(
        `${paths.payments.list}${toQueryString({ ...params })}`,
      ),
  });

export const paymentQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: paymentKeys.detail(uid),
    queryFn: () => api.get<Payment>(paths.payments.payment(uid)),
    retry: retryUnlessNotFound,
  });

export const refundsQueryOptions = (params: RefundListParams) =>
  queryOptions({
    queryKey: paymentKeys.refundList(params),
    queryFn: () =>
      api.get<Paged<Refund>>(
        `${paths.refunds.list}${toQueryString({ ...params })}`,
      ),
  });

export const refundQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: paymentKeys.refund(uid),
    queryFn: () => api.get<Refund>(paths.refunds.refund(uid)),
    retry: retryUnlessNotFound,
  });
