export const PAGE_SIZE = 20;

export interface ListSearch {
  page: number;
  sort?: string;
}

const UUID_PATTERN =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
const PLAIN_DATE_PATTERN = /^\d{4}-\d{2}-\d{2}$/;

export function isUuid(value: string): boolean {
  return UUID_PATTERN.test(value);
}

export function parsePage(value: unknown): number {
  const page = Number(value);
  return Number.isInteger(page) && page >= 0 ? page : 0;
}

export function parseText(value: unknown): string | undefined {
  if (typeof value === "number" && Number.isFinite(value)) return String(value);
  return typeof value === "string" && value !== "" ? value : undefined;
}

export function parseOneOf<T extends string>(
  value: unknown,
  allowed: readonly T[],
): T | undefined {
  return typeof value === "string" &&
    (allowed as readonly string[]).includes(value)
    ? (value as T)
    : undefined;
}

export function parseUuid(value: unknown): string | undefined {
  return typeof value === "string" && isUuid(value) ? value : undefined;
}

export function parsePlainDate(value: unknown): string | undefined {
  return typeof value === "string" && PLAIN_DATE_PATTERN.test(value)
    ? value
    : undefined;
}

export function parseListSearch(search: Record<string, unknown>): ListSearch {
  return {
    page: parsePage(search.page),
    sort: parseText(search.sort),
  };
}
