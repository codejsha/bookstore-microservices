import { expect, test } from "./fixtures";
import { toast } from "./support/locators";

test.use({ authenticated: true });

test.describe("account", () => {
  test("profile tab shows the customer details", async ({ page }) => {
    await page.goto("/account");

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
    await page.goto("/account");
    const firstName = page.getByLabel("First name");
    await expect(firstName).toHaveValue("John");

    await firstName.fill("Jane");
    await page.getByRole("button", { name: "Save changes" }).click();

    await expect(toast(page, "Profile saved")).toBeVisible();

    await page.reload();

    await expect(page.getByLabel("First name")).toHaveValue("Jane");
  });

  test("invalid email is flagged on blur", async ({ page }) => {
    await page.goto("/account");
    const email = page.getByLabel("Email");
    await expect(email).toHaveValue("john.doe@example.com");

    await email.fill("not-an-email");
    await email.blur();

    await expect(page.getByRole("alert")).toHaveText("Invalid email address");
    await expect(email).toHaveAttribute("aria-invalid", "true");
  });

  test("points tab shows the balance", async ({ page }) => {
    await page.goto("/account");

    await page.getByRole("tab", { name: "Points" }).click();

    await expect(page.getByText("Current balance")).toBeVisible();
    await expect(page.getByText("1,240 pt")).toBeVisible();
  });

  test("history tab lists point changes", async ({ page }) => {
    await page.goto("/account");

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
