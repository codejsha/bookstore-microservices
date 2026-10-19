import { expect, test } from "./fixtures";
import {
  dataRows,
  pageSummary,
  rowWithCell,
  selectFilter,
} from "./support/locators";

test.describe("settlements listing", () => {
  test("lists the first page of settlements", async ({ page }) => {
    await page.goto("/settlements");

    await expect(
      page.getByRole("heading", { level: 1, name: "Settlements" }),
    ).toBeVisible();
    await expect(dataRows(page)).toHaveCount(20);
    await expect(pageSummary(page, "1–20 of 24")).toBeVisible();
    const first = dataRows(page).first();
    await expect(first).toContainText("Sep 1, 2026");
    await expect(first).toContainText("USD");
    await expect(first).toContainText("card");
    await expect(first).toContainText(/USD\s27\.00/);
    await expect(first).toContainText(/USD\s8\.22/);
    await expect(first).toContainText("CONFIRMED");
  });

  test("date filter narrows the list to one day", async ({ page }) => {
    await page.goto("/settlements");

    await page.getByLabel("Settlement date").fill("2026-09-02");

    await expect(page).toHaveURL(/date=2026-09-02/);
    await expect(pageSummary(page, "1–3 of 3")).toBeVisible();
    await expect(dataRows(page).first()).toContainText("Sep 2, 2026");
  });

  test("status filter narrows the list", async ({ page }) => {
    await page.goto("/settlements");

    await selectFilter(page, "Status", "OPEN");

    await expect(page).toHaveURL(/status=OPEN/);
    await expect(pageSummary(page, "1–3 of 3")).toBeVisible();
    await expect(dataRows(page).first()).toContainText("Sep 8, 2026");
  });

  test("row click opens the settlement detail", async ({ page, mockApi }) => {
    const settlement = mockApi.data.settlements[0];
    await page.goto("/settlements");

    await page.getByRole("cell", { name: "Sep 1, 2026" }).first().click();

    await expect(page).toHaveURL(`/settlements/${settlement.uid}`);
    await expect(
      page.getByRole("heading", {
        level: 1,
        name: "Sep 1, 2026 · USD · card",
      }),
    ).toBeVisible();
  });
});

test.describe("settlement detail", () => {
  test("shows the totals and detail lines", async ({ page, mockApi }) => {
    const settlement = mockApi.data.settlements[0];
    await page.goto(`/settlements/${settlement.uid}`);

    const totals = page.getByRole("definition");

    await expect(page.getByText(settlement.uid, { exact: true })).toBeVisible();
    await expect(page.getByText("CONFIRMED", { exact: true })).toBeVisible();
    await expect(totals.filter({ hasText: /^USD\s27\.00$/ })).toBeVisible();
    await expect(totals.filter({ hasText: /^USD\s18\.00$/ })).toBeVisible();
    await expect(totals.filter({ hasText: /^USD\s0\.78$/ })).toBeVisible();
    await expect(totals.filter({ hasText: /^USD\s8\.22$/ })).toBeVisible();
    await expect(dataRows(page)).toHaveCount(3);
    await expect(rowWithCell(page, "REFUND")).toContainText(/-USD\s18\.00/);
  });

  test("back link returns to the listing", async ({ page, mockApi }) => {
    await page.goto(`/settlements/${mockApi.data.settlements[0].uid}`);

    await page.getByRole("link", { name: "Back to settlements" }).click();

    await expect(page).toHaveURL(/\/settlements\?page=0$/);
  });

  test("unknown settlement shows the not found state", async ({ page }) => {
    await page.goto("/settlements/bbbbbbbb-bbbb-4bbb-8bbb-999999999999");

    await expect(page.getByText("Settlement not found")).toBeVisible();
  });
});
