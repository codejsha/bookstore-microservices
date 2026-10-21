import { test as base } from "@playwright/test";
import { seedStoredSession } from "./auth";
import { installMockApi, type MockApi } from "./mock-api";

interface Fixtures {
  authenticated: boolean;
  mockApi: MockApi;
}

export const test = base.extend<Fixtures>({
  authenticated: [false, { option: true }],
  context: async ({ context, authenticated }, use) => {
    if (authenticated) {
      await seedStoredSession(context);
    }
    await use(context);
  },
  mockApi: [
    async ({ page }, use) => {
      await use(await installMockApi(page));
    },
    { auto: true },
  ],
});

export { expect } from "@playwright/test";
