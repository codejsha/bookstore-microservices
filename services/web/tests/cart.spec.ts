import { expect, test } from "./fixtures";
import { accountSidebar, cartItem, cartItems, toast } from "./support/locators";
import { serverError } from "./support/requests";

const CART_URL = /\/api\/v1\/cart$/;
const UUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

test.use({ authenticated: true });

test.describe("cart", () => {
  test("lists items with prices and totals", async ({ page }) => {
    await page.goto("/cart");

    await expect(
      page.getByRole("heading", { level: 1, name: "Shopping Cart" }),
    ).toBeVisible();
    await expect(page.getByText("(3 items)")).toBeVisible();
    await expect(cartItems(page)).toHaveCount(2);
    const first = cartItem(page, "Clean Architecture");
    await expect(first).toContainText("$19.99 each");
    await expect(first).toContainText("$19.99");
    const second = cartItem(page, "Designing Data-Intensive Applications");
    await expect(second).toContainText("$41.49 each");
    await expect(second).toContainText("$82.98");
    await expect(page.getByText("Total items: 3")).toBeVisible();
    await expect(page.getByText("$102.97")).toBeVisible();
  });

  test("increasing a quantity updates subtotals and the badge", async ({
    page,
  }) => {
    await page.goto("/cart");
    const item = cartItem(page, "Designing Data-Intensive Applications");

    await item
      .getByRole("button", {
        name: "Increase quantity of Designing Data-Intensive Applications",
      })
      .click();

    await expect(item).toContainText("$124.47");
    await expect(page.getByText("(4 items)")).toBeVisible();
    await expect(page.getByText("$144.46")).toBeVisible();
    await expect(
      accountSidebar(page).getByRole("link", { name: "Cart, 4 items" }),
    ).toBeVisible();
  });

  test("decreasing a single item removes it", async ({ page }) => {
    await page.goto("/cart");

    await cartItem(page, "Clean Architecture")
      .getByRole("button", { name: "Decrease quantity of Clean Architecture" })
      .click();

    await expect(toast(page, "Item removed")).toBeVisible();
    await expect(cartItems(page)).toHaveCount(1);
    await expect(page.getByText("(2 items)")).toBeVisible();
  });

  test("remove action deletes the item", async ({ page }) => {
    await page.goto("/cart");

    await cartItem(page, "Designing Data-Intensive Applications")
      .getByRole("button", {
        name: "Remove Designing Data-Intensive Applications",
      })
      .click();

    await expect(toast(page, "Item removed")).toBeVisible();
    await expect(cartItems(page)).toHaveCount(1);
    await expect(cartItem(page, "Clean Architecture")).toBeVisible();
  });

  test("clearing the cart shows the empty state", async ({ page }) => {
    await page.goto("/cart");

    await page.getByRole("button", { name: "Clear Cart" }).click();

    await expect(toast(page, "Cart cleared")).toBeVisible();
    await expect(page.getByText("Your cart is empty")).toBeVisible();
    await expect(
      accountSidebar(page).getByRole("link", { name: "Cart", exact: true }),
    ).toBeVisible();

    await page.getByRole("link", { name: "Browse Books" }).click();

    await expect(page).toHaveURL("/books");
  });

  test("checkout places an order with an idempotency key", async ({ page }) => {
    await page.goto("/cart");
    const checkoutRequest = page.waitForRequest(
      (request) =>
        request.method() === "POST" &&
        request.url().endsWith("/api/v1/cart/checkout"),
    );

    await page.getByRole("button", { name: "Checkout" }).click();

    const body = (await checkoutRequest).postDataJSON() as {
      currency: string;
      idempotency_key: string;
    };
    expect(body.currency).toBe("USD");
    expect(body.idempotency_key).toMatch(UUID);
    await expect(page).toHaveURL("/orders");
    await expect(toast(page, "Order placed")).toBeVisible();
    await expect(page.getByText("4 orders")).toBeVisible();
    await expect(page.getByRole("link", { name: /ORD-\d{13}/ })).toBeVisible();
  });

  test("empty cart response shows the empty state", async ({ page }) => {
    await page.route(CART_URL, (route) =>
      route.fulfill({
        json: {
          uid: "cart-1",
          user_id: 1,
          items: [],
          total_items: 0,
          total_amount: 0,
        },
      }),
    );
    await page.goto("/cart");

    await expect(page.getByText("Your cart is empty")).toBeVisible();
    await expect(cartItems(page)).toHaveCount(0);
  });

  test("server failure shows an error state", async ({ page }) => {
    await page.route(CART_URL, (route) => route.fulfill(serverError()));
    await page.goto("/cart");

    await expect(page.getByText("Couldn't load your cart")).toBeVisible({
      timeout: 15_000,
    });
    await expect(page.getByRole("button", { name: "Try again" })).toBeVisible();
  });
});
