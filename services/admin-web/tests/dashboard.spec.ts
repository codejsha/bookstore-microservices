import { expect, test } from "./fixtures";
import { serverError } from "./support/requests";

const DASHBOARD_URL = /\/api\/v1\/admin\/dashboard$/;

test.describe("dashboard", () => {
  test("metric cards show counts and flag unavailable sources", async ({
    page,
  }) => {
    await page.goto("/");
    const main = page.getByRole("main");

    await expect(main.getByText("Works", { exact: true })).toBeVisible();
    await expect(main.getByText("1284", { exact: true })).toBeVisible();
    await expect(main.getByText("Orders", { exact: true })).toBeVisible();
    await expect(main.getByText("Unavailable")).toBeVisible();
    await expect(main.getByText("Warehouses", { exact: true })).toBeVisible();
    await expect(main.getByText("3", { exact: true })).toBeVisible();
  });

  test("failed dashboard load shows the page error and recovers on retry", async ({
    page,
  }) => {
    await page.route(DASHBOARD_URL, (route) => route.fulfill(serverError()));
    await page.goto("/");

    await expect(page.getByText("Unable to load this page")).toBeVisible({
      timeout: 15_000,
    });

    await page.unroute(DASHBOARD_URL);
    await page.getByRole("button", { name: "Try again" }).click();

    await expect(
      page.getByRole("heading", { level: 1, name: "Dashboard" }),
    ).toBeVisible();
    await expect(page.getByText("1284", { exact: true })).toBeVisible();
  });
});
