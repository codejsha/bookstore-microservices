import type { ReviewCreateRequestSchema } from "@bookstore/customer-client/model/review-create-request";
import type { ReviewFindResponseSchema } from "@bookstore/customer-client/model/review-find-response";
import type { z } from "zod";

export type Review = Omit<
  z.infer<typeof ReviewFindResponseSchema>,
  "created_at" | "updated_at"
> & {
  created_at: string;
  updated_at?: string;
};

export interface ReviewFindAllResp {
  total: number;
  items: Review[];
}

export type ReviewCreateReq = z.infer<typeof ReviewCreateRequestSchema>;
