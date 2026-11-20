import { expect, test } from "./fixtures";
import { toast } from "./support/locators";
import { localePath } from "./support/paths";

test.use({ authenticated: true });

test.describe("account", () => {
  test("profile tab shows the customer details", async ({ page }) => {
    await page.goto(localePath("/account"));

    await expect(
      page.getByRole("heading", { level: 1, name: "My Account" }),
    ).toBeVisible();
    await expect(page.getByRole("tab", { name: "Profile" })).toHaveAttribute(
      "aria-selected",
      "true",
    );
    await expect(page.getByLabel("First name")).toHaveValue("John");
    await expect(page.getByLabel("Last name")).toHaveValue("Doe");
    await expect(page.getByLabel("Email")).toHaveValue("john.doe@example.com");
    await expect(page.getByLabel("Phone")).toHaveValue("+1-555-0100");
  });

  test("saving the profile persists the change", async ({ page }) => {
    await page.goto(localePath("/account"));
    const firstName = page.getByLabel("First name");
    await expect(firstName).toHaveValue("John");

    await firstName.fill("Jane");
    const [request] = await Promise.all([
      page.waitForRequest(
        (r) => r.method() === "PUT" && /\/api\/v1\/users\/[^/]+$/.test(r.url()),
      ),
      page.getByRole("button", { name: "Save changes" }).click(),
    ]);

    expect(request.postDataJSON()).toEqual({
      first_name: "Jane",
      last_name: "Doe",
      phone: "+1-555-0100",
    });
    await expect(toast(page, "Profile saved")).toBeVisible();

    await page.reload();

    await expect(page.getByLabel("First name")).toHaveValue("Jane");
  });

  test("email is read-only", async ({ page }) => {
    await page.goto(localePath("/account"));
    const email = page.getByLabel("Email");

    await expect(email).toHaveValue("john.doe@example.com");
    await expect(email).not.toBeEditable();
  });

  test("points tab shows the balance", async ({ page }) => {
    await page.goto(localePath("/account"));

    await page.getByRole("tab", { name: "Points" }).click();

    await expect(page.getByText("Current balance")).toBeVisible();
    await expect(page.getByText("1,240 pt")).toBeVisible();
  });

  test("history tab lists point changes", async ({ page }) => {
    await page.goto(localePath("/account"));

    await page.getByRole("tab", { name: "History" }).click();

    const entries = page.getByRole("tabpanel").getByRole("listitem");
    await expect(entries).toHaveCount(6);
    await expect(entries.filter({ hasText: "Welcome bonus" })).toContainText(
      "+100 pt",
    );
    await expect(
      entries.filter({ hasText: "Discount applied to order ORD-2024-0002" }),
    ).toContainText("-200 pt");
  });
});
