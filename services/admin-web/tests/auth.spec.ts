import { expect, test } from "./fixtures";
import { sidebar } from "./support/locators";
import { ME_URL, problem } from "./support/requests";

const SIGN_IN_URL = /\/oauth2\/sign_in/;

test.describe("identity gate", () => {
  test("unauthenticated visitor is sent to sign in with a return path", async ({
    page,
  }) => {
    await page.route(ME_URL, (route) =>
      route.fulfill(problem(401, "Unauthorized")),
    );
    await page.route(SIGN_IN_URL, (route) =>
      route.fulfill({
        contentType: "text/html",
        body: "<title>Sign in</title>",
      }),
    );

    await page.goto("/orders?page=0&status=PAID");

    await page.waitForURL(SIGN_IN_URL);
    expect(new URL(page.url()).searchParams.get("rd")).toBe(
      "/orders?page=0&status=PAID",
    );
  });

  test("forbidden identity lands on the unauthorized page", async ({
    page,
  }) => {
    await page.route(ME_URL, (route) =>
      route.fulfill(problem(403, "Forbidden")),
    );

    await page.goto("/");

    await expect(page).toHaveURL("/unauthorized");
    await expect(
      page.getByRole("heading", { level: 1, name: "Not authorized" }),
    ).toBeVisible();
    await expect(page.getByRole("link", { name: "Sign out" })).toHaveAttribute(
      "href",
      "/oauth2/sign_out",
    );
  });

  test("identity without the staff role lands on the unauthorized page", async ({
    page,
    mockApi,
  }) => {
    await page.route(ME_URL, (route) =>
      route.fulfill({ json: { ...mockApi.data.identity, roles: ["USER"] } }),
    );

    await page.goto("/catalog");

    await expect(page).toHaveURL("/unauthorized");
    await expect(page.getByText("does not have the STAFF role")).toBeVisible();
  });

  test("staff identity opens the console and shows who is signed in", async ({
    page,
  }) => {
    await page.goto("/");

    await expect(
      page.getByRole("heading", { level: 1, name: "Dashboard" }),
    ).toBeVisible();
    await expect(sidebar(page).getByText("admin@example.com")).toBeVisible();
    await expect(
      sidebar(page).getByRole("link", { name: "Sign out" }),
    ).toHaveAttribute("href", "/oauth2/sign_out");
  });

  test("identity is fetched once per session", async ({ page }) => {
    const identityRequests: string[] = [];
    page.on("request", (request) => {
      if (ME_URL.test(request.url())) identityRequests.push(request.url());
    });

    await page.goto("/");
    await sidebar(page).getByRole("link", { name: "Catalog" }).click();
    await expect(
      page.getByRole("heading", { level: 1, name: "Works" }),
    ).toBeVisible();
    await sidebar(page).getByRole("link", { name: "Users" }).click();
    await expect(
      page.getByRole("heading", { level: 1, name: "Users" }),
    ).toBeVisible();

    expect(identityRequests).toHaveLength(1);
  });
});
