import { expect, test } from "./fixtures";
import {
  dataRows,
  pageSummary,
  rowWithCell,
  selectFilter,
} from "./support/locators";

test.describe("orders listing", () => {
  test("lists the first page of orders", async ({ page, mockApi }) => {
    const order = mockApi.data.orders[0];
    await page.goto("/orders");

    await expect(
      page.getByRole("heading", { level: 1, name: "Orders" }),
    ).toBeVisible();
    await expect(dataRows(page)).toHaveCount(20);
    await expect(pageSummary(page, "1–20 of 45")).toBeVisible();
    await expect(rowWithCell(page, order.order_number)).toContainText(
      "PENDING",
    );
    await expect(rowWithCell(page, order.order_number)).toContainText(
      /KRW\s23,790/,
    );
  });

  test("status filter narrows the list", async ({ page }) => {
    await page.goto("/orders");

    await selectFilter(page, "Status", "PENDING");

    await expect(page).toHaveURL(/status=PENDING/);
    await expect(pageSummary(page, "1–8 of 8")).toBeVisible();
    await expect(dataRows(page).last()).toContainText("PENDING");
  });

  test("user uid filter narrows the list", async ({ page, mockApi }) => {
    const user = mockApi.data.users[0];
    await page.goto("/orders");

    await page.getByRole("textbox", { name: "User uid" }).fill(user.uid);

    await expect(page).toHaveURL(new RegExp(`user_uid=${user.uid}`));
    await expect(pageSummary(page, "1–4 of 4")).toBeVisible();
  });

  test("malformed uid marks the input invalid without filtering", async ({
    page,
  }) => {
    await page.goto("/orders");
    const input = page.getByRole("textbox", { name: "User uid" });

    await input.fill("not-a-uuid");

    await expect(input).toHaveAttribute("aria-invalid", "true");
    await expect(page).not.toHaveURL(/user_uid=/);
    await expect(pageSummary(page, "1–20 of 45")).toBeVisible();
  });

  test("user without orders shows the empty message", async ({
    page,
    mockApi,
  }) => {
    await page.goto("/orders");

    await page
      .getByRole("textbox", { name: "User uid" })
      .fill(mockApi.data.users[20].uid);

    await expect(
      page.getByText("No orders match these filters."),
    ).toBeVisible();
  });

  test("user link in a row opens the user instead of the order", async ({
    page,
    mockApi,
  }) => {
    const order = mockApi.data.orders[0];
    await page.goto("/orders");

    await rowWithCell(page, order.order_number)
      .getByRole("link", { name: order.user_uid })
      .click();

    await expect(page).toHaveURL(`/users/${order.user_uid}`);
  });

  test("row click opens the order detail", async ({ page, mockApi }) => {
    const order = mockApi.data.orders[0];
    await page.goto("/orders");

    await page
      .getByRole("cell", { name: order.order_number, exact: true })
      .click();

    await expect(page).toHaveURL(`/orders/${order.uid}`);
    await expect(
      page.getByRole("heading", { level: 1, name: order.order_number }),
    ).toBeVisible();
  });
});

test.describe("order detail", () => {
  test("shows amounts, lines and shipping", async ({ page, mockApi }) => {
    const order = mockApi.data.orders[0];
    await page.goto(`/orders/${order.uid}`);

    await expect(page.getByText(order.uid, { exact: true })).toBeVisible();
    await expect(page.getByText("PENDING", { exact: true })).toBeVisible();
    await expect(page.getByText(/KRW\s2,100/)).toBeVisible();
    await expect(page.getByText(/KRW\s1,890/)).toBeVisible();
    await expect(page.getByText(/KRW\s23,790/)).toBeVisible();
    await expect(dataRows(page)).toHaveCount(1);
    await expect(rowWithCell(page, "Dune")).toContainText("10%");
    await expect(rowWithCell(page, "Dune")).toContainText(/KRW\s21,000/);
    await expect(page.getByText("Ada Lovelace")).toBeVisible();
    await expect(
      page.getByRole("definition").filter({ hasText: /^Seoul$/ }),
    ).toHaveCount(2);
    await expect(page.getByText("express", { exact: true })).toBeVisible();
  });

  test("user link opens the user detail", async ({ page, mockApi }) => {
    const order = mockApi.data.orders[0];
    await page.goto(`/orders/${order.uid}`);

    await page.getByRole("link", { name: order.user_uid }).click();

    await expect(page).toHaveURL(`/users/${order.user_uid}`);
    await expect(
      page.getByRole("heading", { level: 1, name: "Ada Lovelace" }),
    ).toBeVisible();
  });

  test("back link returns to the listing", async ({ page, mockApi }) => {
    await page.goto(`/orders/${mockApi.data.orders[0].uid}`);

    await page.getByRole("link", { name: "Back to orders" }).click();

    await expect(page).toHaveURL(/\/orders\?page=0$/);
  });

  test("unknown order shows the not found state", async ({ page }) => {
    await page.goto("/orders/77777777-7777-4777-8777-999999999999");

    await expect(page.getByText("Order not found")).toBeVisible();
  });
});
