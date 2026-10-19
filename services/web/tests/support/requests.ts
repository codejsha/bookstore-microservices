import type { Request } from "@playwright/test";

export function apiQuery(
  request: Request,
  pathname: string,
): URLSearchParams | undefined {
  const url = new URL(request.url());
  return url.pathname === pathname ? url.searchParams : undefined;
}

export function serverError(title = "Internal Server Error") {
  return {
    status: 500,
    contentType: "application/json",
    body: JSON.stringify({ title, status: 500 }),
  };
}
