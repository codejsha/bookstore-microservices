import { queryOptions } from "@tanstack/react-query";
import { api, retryUnlessNotFound } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import { toQueryString } from "@/shared/api/query-string";
import type { Paged } from "@/shared/api/types";
import type {
  Stock,
  StockListParams,
  Warehouse,
  WarehouseListParams,
} from "./types";

export const inventoryKeys = {
  all: ["inventory"] as const,
  warehouseList: (params: WarehouseListParams) =>
    [...inventoryKeys.all, "warehouse-list", params] as const,
  warehouse: (uid: string) => [...inventoryKeys.all, "warehouse", uid] as const,
  stockList: (params: StockListParams) =>
    [...inventoryKeys.all, "stock-list", params] as const,
};

export const warehousesQueryOptions = (params: WarehouseListParams) =>
  queryOptions({
    queryKey: inventoryKeys.warehouseList(params),
    queryFn: () =>
      api.get<Paged<Warehouse>>(
        `${paths.inventory.warehouses}${toQueryString({ ...params })}`,
      ),
  });

export const warehouseQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: inventoryKeys.warehouse(uid),
    queryFn: () => api.get<Warehouse>(paths.inventory.warehouse(uid)),
    retry: retryUnlessNotFound,
  });

export const stocksQueryOptions = (params: StockListParams) =>
  queryOptions({
    queryKey: inventoryKeys.stockList(params),
    queryFn: () =>
      api.get<Paged<Stock>>(
        `${paths.inventory.stocks}${toQueryString({ ...params })}`,
      ),
  });
