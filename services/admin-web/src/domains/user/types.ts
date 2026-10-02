import type { AdminUserResponseSchema } from "@bookstore/admin-client/model/admin-user-response";
import type { z } from "zod";
import type { ListParams, Wire } from "@/shared/api/types";

export type User = Wire<
  z.infer<typeof AdminUserResponseSchema>,
  "created_at" | "updated_at" | "last_login_at"
>;

export const USER_STATUSES = ["ACTIVE", "SUSPENDED", "DEACTIVATED"] as const;
export type UserStatus = (typeof USER_STATUSES)[number];

export interface UserListParams extends ListParams {
  email?: string;
  name?: string;
  phone?: string;
  status?: UserStatus;
}
