import { expect, test } from "./fixtures";
import { accountSidebar, toast } from "./support/locators";

test.describe("book detail", () => {
  test("renders bibliographic data, editions and reviews", async ({ page }) => {
    await page.goto("/books/work-1");

    await expect(
      page.getByRole("heading", { level: 1, name: "Clean Architecture" }),
    ).toBeVisible();
    await expect(
      page.getByText("by Martin Fowler, Robert C. Martin"),
    ).toBeVisible();
    await expect(page.getByText("First published 2005")).toBeVisible();
    await expect(
      page.getByRole("heading", { level: 2, name: "Description" }),
    ).toBeVisible();
    await expect(
      page.getByRole("heading", { level: 2, name: "About the Authors" }),
    ).toBeVisible();
    await expect(page.getByText("Software architecture expert")).toBeVisible();
    await expect(page.getByText("Clean Code author")).toBeVisible();

    await expect(
      page.getByRole("heading", { level: 2, name: "Editions (2)" }),
    ).toBeVisible();
    await expect(
      page.getByText("Clean Architecture (Hardcover)"),
    ).toBeVisible();
    await expect(page.getByText("$16.49")).toBeVisible();
    await expect(
      page.getByText("Clean Architecture (Paperback)"),
    ).toBeVisible();
    await expect(page.getByText("$17.99")).toBeVisible();

    await expect(page.getByText("Reviews (3)")).toBeVisible();
    await expect(page.getByText("4.7 out of 5")).toBeVisible();
    await expect(
      page.getByText("A must-read for every developer"),
    ).toBeVisible();
  });

  test("back link returns to the listing", async ({ page }) => {
    await page.goto("/books/work-1");

    await page.getByRole("link", { name: "Back to books" }).click();

    await expect(page).toHaveURL("/books");
  });

  test("unknown book shows the page error state", async ({ page }) => {
    await page.goto("/books/does-not-exist");

    await expect(page.getByText("Page not found")).toBeVisible({
      timeout: 15_000,
    });
    await expect(page.getByRole("link", { name: "Go home" })).toBeVisible();
    await expect(page.getByRole("banner")).toHaveCount(1);
  });

  test("wishlist and review submission stay disabled for anonymous visitors", async ({
    page,
  }) => {
    await page.goto("/books/work-1");

    await expect(
      page.getByRole("button", { name: "Add to wishlist" }),
    ).toBeDisabled();

    await page.getByRole("button", { name: "Write a Review" }).click();
    await page
      .getByRole("group", { name: "Your rating" })
      .getByRole("button", { name: "Rate 5 stars" })
      .click();

    await expect(
      page.getByRole("button", { name: "Submit Review" }),
    ).toBeDisabled();
  });
});

test.describe("book detail for signed-in users", () => {
  test.use({ authenticated: true });

  test("adding an edition to the cart updates the badge", async ({ page }) => {
    await page.goto("/books/work-1");
    await expect(
      accountSidebar(page).getByRole("link", { name: "Cart, 3 items" }),
    ).toBeVisible();

    await page.getByRole("button", { name: "Add to Cart" }).first().click();

    await expect(toast(page, "Added to cart")).toContainText(
      "Clean Architecture (Hardcover)",
    );
    await expect(
      accountSidebar(page).getByRole("link", { name: "Cart, 4 items" }),
    ).toBeVisible();
  });

  test("wishlist button toggles membership", async ({ page }) => {
    await page.goto("/books/work-1");
    const add = page.getByRole("button", { name: "Add to wishlist" });
    const remove = page.getByRole("button", { name: "Remove from wishlist" });

    await add.click();

    await expect(toast(page, "Added to wishlist")).toBeVisible();
    await expect(remove).toBeVisible();

    await remove.click();

    await expect(toast(page, "Removed from wishlist")).toBeVisible();
    await expect(add).toBeVisible();
  });

  test("already wished book starts with the remove action", async ({
    page,
  }) => {
    await page.goto("/books/work-2");

    await expect(
      page.getByRole("button", { name: "Remove from wishlist" }),
    ).toBeEnabled();
  });

  test("submitting a review appends it to the list", async ({ page }) => {
    await page.goto("/books/work-1");
    await expect(page.getByText("Reviews (3)")).toBeVisible();

    await page.getByRole("button", { name: "Write a Review" }).click();
    const submit = page.getByRole("button", { name: "Submit Review" });
    await expect(submit).toBeDisabled();

    await page
      .getByRole("group", { name: "Your rating" })
      .getByRole("button", { name: "Rate 5 stars" })
      .click();
    await page.getByLabel("Review title").fill("Great read");
    await page
      .getByLabel("Review content")
      .fill("Clear explanations and practical advice.");
    await submit.click();

    await expect(page.getByText("Reviews (4)")).toBeVisible();
    await expect(page.getByText("Great read", { exact: true })).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Write a Review" }),
    ).toBeVisible();
  });
});
