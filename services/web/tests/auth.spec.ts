import { expect, test } from "./fixtures";

const PROTECTED_ROUTES = ["/cart", "/orders", "/wishlist", "/account"];

test.describe("authentication", () => {
  test("login page offers the Keycloak sign-in", async ({ page }) => {
    await page.goto("/login");
    const main = page.getByRole("main");

    await expect(main.getByText("Login", { exact: true })).toBeVisible();
    await expect(
      main.getByRole("button", { name: "Continue with Keycloak" }),
    ).toBeEnabled();
    await expect(
      main.getByRole("button", { name: "Create an account" }),
    ).toBeEnabled();

    await main.getByRole("link", { name: "Back to home" }).click();

    await expect(page).toHaveURL("/");
  });

  test("signup page links back to login", async ({ page }) => {
    await page.goto("/signup");
    const main = page.getByRole("main");

    await expect(main.getByText("Create account")).toBeVisible();
    await expect(
      main.getByRole("button", { name: "Continue with Keycloak" }),
    ).toBeEnabled();

    await main
      .getByRole("link", { name: "Already have an account? Login" })
      .click();

    await expect(page).toHaveURL("/login");
  });

  for (const route of PROTECTED_ROUTES) {
    test(`${route} redirects anonymous visitors to login`, async ({ page }) => {
      await page.goto(route);

      await expect(page).toHaveURL("/login");
      await expect(
        page.getByRole("button", { name: "Continue with Keycloak" }),
      ).toBeVisible();
    });
  }
});

test.describe("authentication for signed-in users", () => {
  test.use({ authenticated: true });

  for (const route of PROTECTED_ROUTES) {
    test(`${route} stays on the page`, async ({ page }) => {
      await page.goto(route);

      await expect(page).toHaveURL(route);
      await expect(page.getByRole("heading", { level: 1 })).toBeVisible();
    });
  }
});
