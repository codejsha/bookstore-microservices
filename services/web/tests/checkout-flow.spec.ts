import { expect, test } from "./fixtures";
import { accountSidebar, bookItems, cartItem, toast } from "./support/locators";

test.use({ authenticated: true });

test("browse, add to cart, check out and cancel the new order", async ({
  page,
}) => {
  await page.goto("/books");
  await page.getByRole("textbox", { name: "Search by title" }).fill("clean");
  await bookItems(page)
    .filter({ hasText: "Clean Architecture" })
    .getByRole("link")
    .click();
  await expect(
    page.getByRole("heading", { level: 1, name: "Clean Architecture" }),
  ).toBeVisible();

  await page.getByRole("button", { name: "Add to Cart" }).first().click();
  await expect(toast(page, "Added to cart")).toContainText(
    "Clean Architecture (Hardcover)",
  );
  const cartLink = accountSidebar(page).getByRole("link", {
    name: "Cart, 4 items",
  });
  await expect(cartLink).toBeVisible();

  await cartLink.click();
  await expect(page).toHaveURL("/cart");
  await expect(cartItem(page, "Clean Architecture (Hardcover)")).toContainText(
    "$16.49 each",
  );

  await page.getByRole("button", { name: "Checkout" }).click();
  await expect(page).toHaveURL("/orders");
  await expect(toast(page, "Order placed")).toBeVisible();
  const newOrder = page.getByRole("link", { name: /ORD-\d{13}/ });
  await expect(newOrder).toContainText("PENDING");
  await expect(newOrder).toContainText("3 items");

  await newOrder.click();
  await expect(page).toHaveURL(/\/orders\/ord-\d{13}$/);
  await expect(
    page
      .getByRole("main")
      .getByText("Clean Architecture (Hardcover)", { exact: true }),
  ).toBeVisible();

  await page.getByRole("button", { name: "Cancel Order" }).click();
  await expect(toast(page, "Order cancelled")).toBeVisible();
  await expect(
    page.getByRole("main").getByText("CANCELLED", { exact: true }),
  ).toBeVisible();
  await expect(
    accountSidebar(page).getByRole("link", { name: "Cart", exact: true }),
  ).toBeVisible();
});
