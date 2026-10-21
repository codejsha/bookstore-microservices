import { expect, test } from "./fixtures";
import { sidebar } from "./support/locators";

const SECTIONS = [
  { label: "Dashboard", path: "/", heading: "Dashboard" },
  { label: "Catalog", path: "/catalog", heading: "Works" },
  { label: "Users", path: "/users", heading: "Users" },
  { label: "Orders", path: "/orders", heading: "Orders" },
  { label: "Payments", path: "/payments", heading: "Payments" },
  { label: "Refunds", path: "/refunds", heading: "Refunds" },
  { label: "Settlements", path: "/settlements", heading: "Settlements" },
  {
    label: "Warehouses",
    path: "/inventory/warehouses",
    heading: "Warehouses",
  },
  { label: "Stock", path: "/inventory/stocks", heading: "Stock" },
  { label: "Risk", path: "/risk", heading: "Risk" },
];

test.describe("console shell", () => {
  for (const section of SECTIONS) {
    test(`${section.label} link opens its section and becomes current`, async ({
      page,
    }) => {
      await page.goto(section.label === "Dashboard" ? "/catalog" : "/");
      const link = sidebar(page).getByRole("link", {
        name: section.label,
        exact: true,
      });

      await link.click();

      await expect(page).toHaveURL(new RegExp(`^[^?]*${section.path}(\\?|$)`));
      await expect(
        page.getByRole("heading", { level: 1, name: section.heading }),
      ).toBeVisible();
      await expect(link).toHaveAttribute("aria-current", "page");
    });
  }

  test("only the open section is marked current", async ({ page }) => {
    await page.goto("/orders");

    await expect(
      sidebar(page).getByRole("link", { name: "Orders" }),
    ).toHaveAttribute("aria-current", "page");
    await expect(
      sidebar(page).getByRole("link", { name: "Dashboard" }),
    ).not.toHaveAttribute("aria-current", "page");
  });

  test("unknown route shows the console not found page", async ({ page }) => {
    await page.goto("/no-such-section");

    await expect(page.getByText("Page not found")).toBeVisible();
    await expect(
      page.getByText("That page does not exist in the admin console."),
    ).toBeVisible();
  });
});
