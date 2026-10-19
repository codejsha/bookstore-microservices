import type { Page, Route } from "@playwright/test";
import { getResponse } from "msw";
import { createMockData, type MockData } from "../../src/mocks/data";
import { createHandlers } from "../../src/mocks/handlers";

export interface MockApi {
  data: MockData;
}

const API_URL = /\/api\/v1\//;
const TELEMETRY_URL = /\/(otlp|faro)\//;

async function toFetchRequest(route: Route): Promise<Request> {
  const request = route.request();
  const method = request.method();
  const headers = Object.fromEntries(
    Object.entries(await request.allHeaders()).filter(
      ([name]) => !name.startsWith(":"),
    ),
  );
  const body =
    method === "GET" || method === "HEAD"
      ? undefined
      : (request.postData() ?? undefined);
  return new Request(request.url(), { method, headers, body });
}

export async function installMockApi(page: Page): Promise<MockApi> {
  const data = createMockData();
  const handlers = createHandlers(data);

  await page.route(TELEMETRY_URL, (route) => route.fulfill({ status: 204 }));
  await page.route(API_URL, async (route) => {
    const request = await toFetchRequest(route);
    const response = await getResponse(handlers, request, {
      baseUrl: new URL(request.url).origin,
    });
    if (!response) {
      await route.fallback();
      return;
    }
    await route.fulfill({
      status: response.status,
      headers: Object.fromEntries(response.headers),
      body: Buffer.from(await response.arrayBuffer()),
    });
  });

  return { data };
}
