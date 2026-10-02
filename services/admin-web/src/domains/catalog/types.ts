import type { AdminAuthorItemSchema } from "@bookstore/admin-client/model/admin-author-item";
import type { AdminSubjectItemSchema } from "@bookstore/admin-client/model/admin-subject-item";
import type { AdminWorkCreateRequestSchema } from "@bookstore/admin-client/model/admin-work-create-request";
import type { AdminWorkResponseSchema } from "@bookstore/admin-client/model/admin-work-response";
import type { AdminWorkUpdateRequestSchema } from "@bookstore/admin-client/model/admin-work-update-request";
import type { z } from "zod";
import type { ListParams } from "@/shared/api/types";

export type { Paged } from "@/shared/api/types";

export type Work = Omit<
  z.infer<typeof AdminWorkResponseSchema>,
  "created_at" | "updated_at"
> & {
  created_at: string;
  updated_at?: string;
};

export type Author = z.infer<typeof AdminAuthorItemSchema>;
export type Subject = z.infer<typeof AdminSubjectItemSchema>;
export type WorkCreateRequest = z.infer<typeof AdminWorkCreateRequestSchema>;
export type WorkUpdateRequest = z.infer<typeof AdminWorkUpdateRequestSchema>;

export interface WorkListParams extends ListParams {
  title?: string;
}
