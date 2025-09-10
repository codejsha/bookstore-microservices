import type { AuthorCreateRequestSchema } from "@bookstore/catalog-client/model/author-create-request";
import type { AuthorFindAllResponseSchema } from "@bookstore/catalog-client/model/author-find-all-response";
import type { AuthorItemSchema } from "@bookstore/catalog-client/model/author-item";
import type { AuthorUpdateRequestSchema } from "@bookstore/catalog-client/model/author-update-request";
import type { EditionFindResponseSchema } from "@bookstore/catalog-client/model/edition-find-response";
import type { EditionItemSchema } from "@bookstore/catalog-client/model/edition-item";
import type { PublisherCreateRequestSchema } from "@bookstore/catalog-client/model/publisher-create-request";
import type { PublisherFindAllResponseSchema } from "@bookstore/catalog-client/model/publisher-find-all-response";
import type { PublisherItemSchema } from "@bookstore/catalog-client/model/publisher-item";
import type { PublisherUpdateRequestSchema } from "@bookstore/catalog-client/model/publisher-update-request";
import type { SubjectCreateRequestSchema } from "@bookstore/catalog-client/model/subject-create-request";
import type { SubjectFindAllResponseSchema } from "@bookstore/catalog-client/model/subject-find-all-response";
import type { SubjectItemSchema } from "@bookstore/catalog-client/model/subject-item";
import type { SubjectUpdateRequestSchema } from "@bookstore/catalog-client/model/subject-update-request";
import type { WorkCreateRequestSchema } from "@bookstore/catalog-client/model/work-create-request";
import type { WorkFindResponseSchema } from "@bookstore/catalog-client/model/work-find-response";
import type { WorkItemSchema } from "@bookstore/catalog-client/model/work-item";
import type { WorkSearchItemSchema } from "@bookstore/catalog-client/model/work-search-item";
import type { WorkUpdateRequestSchema } from "@bookstore/catalog-client/model/work-update-request";
import type { z } from "zod";

type WithStringDates<T> = Omit<T, "created_at" | "updated_at"> & {
  created_at: string;
  updated_at?: string;
};

// ── Authors ──────────────────────────────────────────────────────────────
export type Author = z.infer<typeof AuthorItemSchema>;
export type AuthorCreateReq = z.infer<typeof AuthorCreateRequestSchema>;
export type AuthorUpdateReq = z.infer<typeof AuthorUpdateRequestSchema>;
export type AuthorFindAllResp = z.infer<typeof AuthorFindAllResponseSchema>;

// ── Subjects ─────────────────────────────────────────────────────────────
export type Subject = z.infer<typeof SubjectItemSchema>;
export type SubjectCreateReq = z.infer<typeof SubjectCreateRequestSchema>;
export type SubjectUpdateReq = z.infer<typeof SubjectUpdateRequestSchema>;
export type SubjectFindAllResp = z.infer<typeof SubjectFindAllResponseSchema>;

// ── Publishers ───────────────────────────────────────────────────────────
export type Publisher = z.infer<typeof PublisherItemSchema>;
export type PublisherCreateReq = z.infer<typeof PublisherCreateRequestSchema>;
export type PublisherUpdateReq = z.infer<typeof PublisherUpdateRequestSchema>;
export type PublisherFindAllResp = z.infer<
  typeof PublisherFindAllResponseSchema
>;

// ── Works ────────────────────────────────────────────────────────────────
export type WorkItem = WithStringDates<z.infer<typeof WorkItemSchema>>;
export type Work = WithStringDates<z.infer<typeof WorkFindResponseSchema>>;
export type WorkSearchItem = WithStringDates<
  z.infer<typeof WorkSearchItemSchema>
>;
export type WorkCreateReq = z.infer<typeof WorkCreateRequestSchema>;
export type WorkUpdateReq = z.infer<typeof WorkUpdateRequestSchema>;

export interface WorkFindAllResp {
  total: number;
  items: WorkItem[];
}

export interface WorkSearchResp {
  total: number;
  items: WorkSearchItem[];
}

// ── Editions ─────────────────────────────────────────────────────────────
export type Edition = WithStringDates<z.infer<typeof EditionItemSchema>>;
export type EditionFindResponse = WithStringDates<
  z.infer<typeof EditionFindResponseSchema>
>;

export interface EditionFindAllResp {
  total: number;
  items: Edition[];
}

// ── Pagination + query params (domain-specific, not generated) ─────────────
export interface PaginationParams {
  size?: number;
  page?: number;
  sort?: string;
}

export interface WorkFindAllParams extends PaginationParams {
  title?: string;
  authorUid?: string;
  subjectUid?: string;
  olKey?: string;
}

export interface WorkSearchParams extends PaginationParams {
  q?: string;
  authorUid?: string;
  subjectUid?: string;
}

export interface EditionFindAllParams extends PaginationParams {
  title?: string;
  isbn?: string;
  workUid?: string;
  publisherUid?: string;
  language?: string;
  olKey?: string;
}

export interface AuthorFindAllParams extends PaginationParams {
  name?: string;
}

export interface PublisherFindAllParams extends PaginationParams {
  name?: string;
}

export interface SubjectFindAllParams extends PaginationParams {
  name?: string;
}
