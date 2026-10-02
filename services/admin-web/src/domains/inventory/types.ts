import type { AdminStockResponseSchema } from "@bookstore/admin-client/model/admin-stock-response";
import type { AdminStockWarehouseItemSchema } from "@bookstore/admin-client/model/admin-stock-warehouse-item";
import type { AdminWarehouseResponseSchema } from "@bookstore/admin-client/model/admin-warehouse-response";
import type { z } from "zod";
import type { ListParams, Wire } from "@/shared/api/types";

export type Warehouse = Wire<
  z.infer<typeof AdminWarehouseResponseSchema>,
  "created_at" | "updated_at"
>;

export type Stock = z.infer<typeof AdminStockResponseSchema>;
export type StockAtWarehouse = z.infer<typeof AdminStockWarehouseItemSchema>;

export interface WarehouseListParams extends ListParams {
  name?: string;
}

export interface StockListParams extends ListParams {
  edition_uid?: string;
  warehouse_uid?: string;
}
