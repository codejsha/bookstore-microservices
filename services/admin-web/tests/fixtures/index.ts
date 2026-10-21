import { test as base } from "@playwright/test";
import { installMockApi, type MockApi } from "./mock-api";

interface Fixtures {
  mockApi: MockApi;
}

export const test = base.extend<Fixtures>({
  mockApi: [
    async ({ page }, use) => {
      await use(await installMockApi(page));
    },
    { auto: true },
  ],
});

export { expect } from "@playwright/test";
