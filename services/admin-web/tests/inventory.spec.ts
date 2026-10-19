import { expect, test } from "./fixtures";
import { dataRows, pageSummary, rowWithCell } from "./support/locators";

test.describe("warehouses", () => {
  test("lists every warehouse", async ({ page }) => {
    await page.goto("/inventory/warehouses");

    await expect(
      page.getByRole("heading", { level: 1, name: "Warehouses" }),
    ).toBeVisible();
    await expect(dataRows(page)).toHaveCount(3);
    await expect(pageSummary(page, "1–3 of 3")).toBeVisible();
    await expect(rowWithCell(page, "Incheon Central")).toContainText("120,000");
    await expect(rowWithCell(page, "Oakland Overflow")).toContainText("—");
  });

  test("name filter narrows the list", async ({ page }) => {
    await page.goto("/inventory/warehouses");

    await page.getByRole("textbox", { name: "Name" }).fill("busan");

    await expect(page).toHaveURL(/name=busan/);
    await expect(dataRows(page)).toHaveCount(1);
    await expect(rowWithCell(page, "Busan Port")).toBeVisible();
  });

  test("row click opens the warehouse detail", async ({ page, mockApi }) => {
    const warehouse = mockApi.data.warehouses[1];
    await page.goto("/inventory/warehouses");

    await page.getByRole("cell", { name: warehouse.name, exact: true }).click();

    await expect(page).toHaveURL(`/inventory/warehouses/${warehouse.uid}`);
    await expect(
      page.getByRole("heading", { level: 1, name: "Busan Port" }),
    ).toBeVisible();
    await expect(page.getByText("80,000")).toBeVisible();
  });

  test("view stock link opens the stock held at the warehouse", async ({
    page,
    mockApi,
  }) => {
    const warehouse = mockApi.data.warehouses[1];
    await page.goto(`/inventory/warehouses/${warehouse.uid}`);

    await page.getByRole("link", { name: "View stock" }).click();

    await expect(page).toHaveURL(new RegExp(`warehouse_uid=${warehouse.uid}`));
    await expect(
      page.getByRole("textbox", { name: "Warehouse uid" }),
    ).toHaveValue(warehouse.uid);
    await expect(pageSummary(page, "1–20 of 28")).toBeVisible();
    await expect(dataRows(page).first()).toContainText("Busan Port");
    await expect(dataRows(page).first()).not.toContainText("Incheon Central");
  });

  test("unknown warehouse shows the not found state", async ({ page }) => {
    await page.goto(
      "/inventory/warehouses/cccccccc-cccc-4ccc-8ccc-999999999999",
    );

    await expect(page.getByText("Warehouse not found")).toBeVisible();
  });
});

test.describe("stock", () => {
  test("lists stock with per-warehouse quantities", async ({
    page,
    mockApi,
  }) => {
    const stock = mockApi.data.stocks[0];
    await page.goto("/inventory/stocks");

    await expect(
      page.getByRole("heading", { level: 1, name: "Stock" }),
    ).toBeVisible();
    await expect(pageSummary(page, "1–20 of 28")).toBeVisible();
    const first = rowWithCell(page, stock.edition_uid);
    await expect(first).toContainText("45");
    await expect(first).toContainText("Incheon Central 17");
    await expect(first).toContainText("Busan Port 28");
  });

  test("edition filter narrows the list", async ({ page, mockApi }) => {
    const stock = mockApi.data.stocks[1];
    await page.goto("/inventory/stocks");

    await page
      .getByRole("textbox", { name: "Edition uid" })
      .fill(stock.edition_uid);

    await expect(page).toHaveURL(
      new RegExp(`edition_uid=${stock.edition_uid}`),
    );
    await expect(dataRows(page)).toHaveCount(1);
    await expect(rowWithCell(page, stock.edition_uid)).toBeVisible();
  });

  test("warehouse link in a row opens the warehouse detail", async ({
    page,
  }) => {
    await page.goto("/inventory/stocks");

    await dataRows(page)
      .first()
      .getByRole("link", { name: "Incheon Central" })
      .click();

    await expect(page).toHaveURL(/\/inventory\/warehouses\//);
    await expect(
      page.getByRole("heading", { level: 1, name: "Incheon Central" }),
    ).toBeVisible();
  });

  test("filter without matches shows the empty message", async ({ page }) => {
    await page.goto("/inventory/stocks");

    await page
      .getByRole("textbox", { name: "Edition uid" })
      .fill("eeeeeeee-eeee-4eee-8eee-999999999999");

    await expect(
      page.getByText("No stock matches these filters."),
    ).toBeVisible();
  });
});
