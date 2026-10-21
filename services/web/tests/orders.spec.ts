import { expect, test } from "./fixtures";
import { toast } from "./support/locators";
import { buildOrders } from "./support/orders";

const ORDERS_URL = /\/api\/v1\/orders(\?|$)/;

test.use({ authenticated: true });

test.describe("orders listing", () => {
  test("lists orders with status and totals", async ({ page }) => {
    await page.goto("/orders");

    await expect(
      page.getByRole("heading", { level: 1, name: "Orders" }),
    ).toBeVisible();
    await expect(page.getByText("3 orders")).toBeVisible();
    const pending = page.getByRole("link", { name: /ORD-2024-0001/ });
    await expect(pending).toContainText("PENDING");
    await expect(pending).toContainText("2 items");
    await expect(pending).toContainText("USD 56.11");
    await expect(
      page.getByRole("link", { name: /ORD-2024-0002/ }),
    ).toContainText("DELIVERED");
    await expect(
      page.getByRole("link", { name: /ORD-2024-0003/ }),
    ).toContainText("CANCELLED");
  });

  test("status filter narrows the list", async ({ page }) => {
    await page.goto("/orders");
    const pending = page.getByRole("button", { name: "Pending" });

    await pending.click();

    await expect(pending).toHaveAttribute("aria-pressed", "true");
    await expect(page.getByText("1 order", { exact: true })).toBeVisible();
    await expect(
      page.getByRole("link", { name: /ORD-2024-0001/ }),
    ).toBeVisible();
    await expect(page.getByRole("link", { name: /ORD-2024-0002/ })).toHaveCount(
      0,
    );

    await page.getByRole("button", { name: "Paid" }).click();

    await expect(page.getByText("No orders found")).toBeVisible();
    await expect(page.getByText("Try a different filter.")).toBeVisible();
  });

  test("pagination moves between pages", async ({ page, mockApi }) => {
    const orders = buildOrders(25, mockApi.data.orders[0]);
    await page.route(ORDERS_URL, (route) => {
      const url = new URL(route.request().url());
      const size = Number(url.searchParams.get("size")) || 10;
      const pageIndex = Number(url.searchParams.get("page")) || 0;
      return route.fulfill({
        json: {
          total: orders.length,
          items: orders.slice(pageIndex * size, (pageIndex + 1) * size),
        },
      });
    });
    await page.goto("/orders");
    const previous = page.getByRole("button", { name: "Previous page" });
    const next = page.getByRole("button", { name: "Next page" });

    await expect(page.getByText("25 orders")).toBeVisible();
    await expect(page.getByText("Page 1 of 3")).toBeVisible();
    await expect(previous).toBeDisabled();
    await expect(
      page.getByRole("link", { name: /ORD-GEN-0001/ }),
    ).toBeVisible();

    await next.click();

    await expect(page.getByText("Page 2 of 3")).toBeVisible();
    await expect(previous).toBeEnabled();
    await expect(
      page.getByRole("link", { name: /ORD-GEN-0011/ }),
    ).toBeVisible();

    await next.click();

    await expect(page.getByText("Page 3 of 3")).toBeVisible();
    await expect(next).toBeDisabled();
  });

  test("order card opens the detail page", async ({ page }) => {
    await page.goto("/orders");

    await page.getByRole("link", { name: /ORD-2024-0002/ }).click();

    await expect(page).toHaveURL("/orders/ord-2");
    await expect(
      page.getByRole("heading", { level: 1, name: "ORD-2024-0002" }),
    ).toBeVisible();
  });
});

test.describe("order detail", () => {
  test("pending order shows items, summary, shipping and can be cancelled", async ({
    page,
  }) => {
    await page.goto("/orders/ord-1");
    const main = page.getByRole("main");

    await expect(
      page.getByRole("heading", { level: 1, name: "ORD-2024-0001" }),
    ).toBeVisible();
    await expect(main.getByText("PENDING", { exact: true })).toBeVisible();
    await expect(page.getByText("Items (2)")).toBeVisible();
    await expect(
      page.getByText("Clean Architecture", { exact: true }),
    ).toBeVisible();
    await expect(
      page.getByText("Domain-Driven Design", { exact: true }),
    ).toBeVisible();
    await expect(page.getByText("USD 56.11")).toBeVisible();
    await expect(page.getByText("John Doe")).toBeVisible();
    await expect(page.getByText("Method: Standard")).toBeVisible();

    await page.getByRole("button", { name: "Cancel Order" }).click();

    await expect(toast(page, "Order cancelled")).toBeVisible();
    await expect(main.getByText("CANCELLED", { exact: true })).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Cancel Order" }),
    ).toHaveCount(0);
  });

  test("delivered order shows the discount and no cancel action", async ({
    page,
  }) => {
    await page.goto("/orders/ord-2");

    await expect(page.getByText("-USD 5.00")).toBeVisible();
    await expect(page.getByText("USD 22.88")).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Cancel Order" }),
    ).toHaveCount(0);
    await expect(page.getByText(/^Method: /)).toHaveCount(0);
  });

  test("back link returns to the listing", async ({ page }) => {
    await page.goto("/orders/ord-1");

    await page.getByRole("link", { name: "Back to orders" }).click();

    await expect(page).toHaveURL("/orders");
  });

  test("unknown order shows the page error state", async ({ page }) => {
    await page.goto("/orders/does-not-exist");

    await expect(page.getByText("Page not found")).toBeVisible({
      timeout: 15_000,
    });
    await expect(page.getByRole("link", { name: "Go home" })).toBeVisible();
    await expect(page.getByRole("banner")).toHaveCount(1);
  });
});
