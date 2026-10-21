import { expect, test } from "./fixtures";

test.describe("home", () => {
  test("browse books call to action opens the listing", async ({ page }) => {
    await page.goto("/");

    await page.getByRole("link", { name: "Browse Books" }).click();

    await expect(page).toHaveURL("/books");
  });

  test("create account call to action opens signup", async ({ page }) => {
    await page.goto("/");

    await page.getByRole("link", { name: "Create Account" }).click();

    await expect(page).toHaveURL("/signup");
  });

  test("trending section lists the newest four books", async ({ page }) => {
    await page.goto("/");
    const trending = page
      .getByRole("list", { name: "Trending books" })
      .getByRole("listitem");

    await expect(
      page.getByRole("heading", { level: 2, name: "Trending Now" }),
    ).toBeVisible();
    await expect(trending).toHaveCount(4);
    await expect(trending.first()).toContainText("Clean Architecture");

    await page.getByRole("link", { name: "View all" }).click();

    await expect(page).toHaveURL("/books");
  });
});
