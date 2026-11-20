import { expect, test } from "./fixtures";
import { toast, wishlistItems } from "./support/locators";
import { localePath } from "./support/paths";

const WISHLIST_URL = /\/api\/v1\/customers\/[^/]+\/wishlist$/;

test.use({ authenticated: true });

test.describe("wishlist", () => {
  test("lists the wished books", async ({ page }) => {
    await page.goto(localePath("/wishlist"));

    await expect(
      page.getByRole("heading", { level: 1, name: "My Wishlist" }),
    ).toBeVisible();
    await expect(page.getByText("(4 books)")).toBeVisible();
    await expect(wishlistItems(page)).toHaveCount(4);
    for (const title of [
      "Domain-Driven Design",
      "Refactoring",
      "Release It!",
      "System Design Interview",
    ]) {
      await expect(
        wishlistItems(page).filter({ hasText: title }),
      ).toBeVisible();
    }
  });

  test("removing a book updates the list", async ({ page }) => {
    await page.goto(localePath("/wishlist"));

    await page
      .getByRole("button", { name: "Remove Refactoring from wishlist" })
      .click();

    await expect(toast(page, "Removed from wishlist")).toBeVisible();
    await expect(wishlistItems(page)).toHaveCount(3);
    await expect(page.getByText("(3 books)")).toBeVisible();
    await expect(
      wishlistItems(page).filter({ hasText: "Refactoring" }),
    ).toHaveCount(0);
  });

  test("book title opens the detail page", async ({ page }) => {
    await page.goto(localePath("/wishlist"));

    await page.getByRole("link", { name: "Domain-Driven Design" }).click();

    await expect(page).toHaveURL(localePath("/books/work-2"));
  });

  test("empty wishlist shows the empty state", async ({ page }) => {
    await page.route(WISHLIST_URL, (route) =>
      route.fulfill({ json: { book_uids: [] } }),
    );
    await page.goto(localePath("/wishlist"));

    await expect(page.getByText("Your wishlist is empty")).toBeVisible();

    await page.getByRole("link", { name: "Browse Books" }).click();

    await expect(page).toHaveURL(localePath("/books"));
  });
});
