import { expect, test } from "./fixtures";
import { accountSidebar, browseSidebar } from "./support/locators";
import { localePath } from "./support/paths";

test.describe("global layout", () => {
  test("header logo link returns home", async ({ page }) => {
    await page.goto(localePath("/books"));

    await page
      .getByRole("banner")
      .getByRole("link", { name: "Bookstore" })
      .click();

    await expect(page).toHaveURL(localePath("/"));
    await expect(
      page.getByRole("heading", {
        level: 1,
        name: "Discover Your Next Great Read",
      }),
    ).toBeVisible();
  });

  test("browse sidebar navigates to the books listing", async ({ page }) => {
    await page.goto(localePath("/"));

    await browseSidebar(page).getByRole("link", { name: "Books" }).click();

    await expect(page).toHaveURL(localePath("/books"));
    await expect(
      page.getByRole("heading", { level: 1, name: "Books" }),
    ).toBeVisible();
  });

  test("footer navigation links point at their routes", async ({ page }) => {
    await page.goto(localePath("/"));
    const footer = page.getByRole("contentinfo");
    const browse = footer.getByRole("navigation", { name: "Browse" });
    const account = footer.getByRole("navigation", { name: "Account" });

    await expect(browse.getByRole("link", { name: "Books" })).toHaveAttribute(
      "href",
      localePath("/books"),
    );
    await expect(
      browse.getByRole("link", { name: "Wishlist" }),
    ).toHaveAttribute("href", localePath("/wishlist"));
    await expect(browse.getByRole("link", { name: "Cart" })).toHaveAttribute(
      "href",
      localePath("/cart"),
    );
    await expect(
      account.getByRole("link", { name: "My account" }),
    ).toHaveAttribute("href", localePath("/account"));
    await expect(account.getByRole("link", { name: "Orders" })).toHaveAttribute(
      "href",
      localePath("/orders"),
    );
    await expect(account.getByRole("link", { name: "Login" })).toHaveAttribute(
      "href",
      localePath("/login"),
    );
    await expect(
      account.getByRole("link", { name: "Sign up" }),
    ).toHaveAttribute("href", localePath("/signup"));
  });

  test("theme toggle switches dark mode and persists across reloads", async ({
    page,
  }) => {
    await page.emulateMedia({ colorScheme: "light" });
    await page.goto(localePath("/"));
    const html = page.locator("html");
    await expect(html).not.toHaveClass(/dark/);

    await page.getByRole("button", { name: "Switch to dark theme" }).click();

    await expect(html).toHaveClass(/dark/);
    await expect(
      page.getByRole("button", { name: "Switch to light theme" }),
    ).toBeVisible();

    await page.reload();

    await expect(html).toHaveClass(/dark/);
  });

  test("unknown route shows the not found page", async ({ page }) => {
    await page.goto(localePath("/no-such-page"));

    await expect(page.getByText("Page not found")).toBeVisible();
    await expect(page.getByRole("banner")).toHaveCount(1);

    await page.getByRole("link", { name: "Go home" }).click();

    await expect(page).toHaveURL(localePath("/"));
  });

  test("cart link carries no badge for anonymous visitors", async ({
    page,
  }) => {
    await page.goto(localePath("/"));

    await expect(
      accountSidebar(page).getByRole("link", { name: "Cart", exact: true }),
    ).toBeVisible();
  });
});

test.describe("global layout for signed-in users", () => {
  test.use({ authenticated: true });

  test("cart link announces the item count", async ({ page }) => {
    await page.goto(localePath("/"));

    await expect(
      accountSidebar(page).getByRole("link", { name: "Cart, 3 items" }),
    ).toBeVisible();
  });
});
