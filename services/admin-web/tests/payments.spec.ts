import { expect, test } from "./fixtures";
import {
  dataRows,
  pageSummary,
  rowWithCell,
  selectFilter,
} from "./support/locators";

test.describe("payments listing", () => {
  test("lists the first page of payments", async ({ page, mockApi }) => {
    const payment = mockApi.data.payments[0];
    await page.goto("/payments");

    await expect(
      page.getByRole("heading", { level: 1, name: "Payments" }),
    ).toBeVisible();
    await expect(dataRows(page)).toHaveCount(20);
    await expect(pageSummary(page, "1–20 of 30")).toBeVisible();
    const first = rowWithCell(page, payment.payment_id);
    await expect(first).toContainText("requires_payment_method");
    await expect(first).toContainText(/KRW\s15,000/);
    await expect(first).toContainText("stripe");
  });

  test("status filter narrows the list", async ({ page }) => {
    await page.goto("/payments");

    await selectFilter(page, "Status", "succeeded");

    await expect(page).toHaveURL(/status=succeeded/);
    await expect(pageSummary(page, "1–3 of 3")).toBeVisible();
  });

  test("connector filter narrows the list", async ({ page }) => {
    await page.goto("/payments");

    await page.getByRole("textbox", { name: "Connector" }).fill("adyen");

    await expect(pageSummary(page, "1–15 of 15")).toBeVisible();
    await expect(dataRows(page).first()).toContainText("adyen");
  });

  test("customer filter narrows the list", async ({ page }) => {
    await page.goto("/payments");

    await page.getByRole("textbox", { name: "Customer ID" }).fill("cus_0001");

    await expect(pageSummary(page, "1–3 of 3")).toBeVisible();
  });

  test("row click opens the payment detail", async ({ page, mockApi }) => {
    const payment = mockApi.data.payments[0];
    await page.goto("/payments");

    await page
      .getByRole("cell", { name: payment.payment_id, exact: true })
      .click();

    await expect(page).toHaveURL(`/payments/${payment.uid}`);
    await expect(
      page.getByRole("heading", { level: 1, name: payment.payment_id }),
    ).toBeVisible();
  });
});

test.describe("payment detail", () => {
  test("shows the payment fields", async ({ page, mockApi }) => {
    const payment = mockApi.data.payments[0];
    await page.goto(`/payments/${payment.uid}`);

    await expect(page.getByText(payment.uid, { exact: true })).toBeVisible();
    await expect(
      page.getByText("requires_payment_method", { exact: true }),
    ).toBeVisible();
    await expect(page.getByText(/KRW\s15,000/)).toBeVisible();
    await expect(page.getByText("card", { exact: true })).toBeVisible();
    await expect(page.getByText("stripe", { exact: true })).toBeVisible();
  });

  test("view refunds link opens the refunds for the payment", async ({
    page,
    mockApi,
  }) => {
    const payment = mockApi.data.payments[0];
    await page.goto(`/payments/${payment.uid}`);

    await page.getByRole("link", { name: "View refunds" }).click();

    await expect(page).toHaveURL(/\/refunds\?/);
    await expect(page.getByRole("textbox", { name: "Payment ID" })).toHaveValue(
      payment.payment_id,
    );
    await expect(pageSummary(page, "1–1 of 1")).toBeVisible();
    await expect(dataRows(page).first()).toContainText(payment.payment_id);
  });

  test("customer link filters the payments by customer", async ({
    page,
    mockApi,
  }) => {
    const payment = mockApi.data.payments[0];
    await page.goto(`/payments/${payment.uid}`);

    await page.getByRole("link", { name: "cus_0000" }).click();

    await expect(page).toHaveURL(/customer_id=cus_0000/);
    await expect(
      page.getByRole("textbox", { name: "Customer ID" }),
    ).toHaveValue("cus_0000");
    await expect(pageSummary(page, "1–4 of 4")).toBeVisible();
  });

  test("unknown payment shows the not found state", async ({ page }) => {
    await page.goto("/payments/88888888-8888-4888-8888-999999999999");

    await expect(page.getByText("Payment not found")).toBeVisible();
  });
});

test.describe("refunds", () => {
  test("lists every refund on one page", async ({ page, mockApi }) => {
    const refund = mockApi.data.refunds[0];
    await page.goto("/refunds");

    await expect(
      page.getByRole("heading", { level: 1, name: "Refunds" }),
    ).toBeVisible();
    await expect(pageSummary(page, "1–14 of 14")).toBeVisible();
    await expect(page.getByRole("button", { name: "Next" })).toBeDisabled();
    const first = rowWithCell(page, refund.refund_id);
    await expect(first).toContainText(refund.payment_id);
    await expect(first).toContainText("pending");
    await expect(first).toContainText("scheduled");
  });

  test("status filter narrows the list", async ({ page }) => {
    await page.goto("/refunds");

    await selectFilter(page, "Status", "pending");

    await expect(page).toHaveURL(/status=pending/);
    await expect(pageSummary(page, "1–4 of 4")).toBeVisible();
  });

  test("row click opens the refund detail with its payment link", async ({
    page,
    mockApi,
  }) => {
    const refund = mockApi.data.refunds[0];
    await page.goto("/refunds");

    await page
      .getByRole("cell", { name: refund.refund_id, exact: true })
      .click();

    await expect(page).toHaveURL(`/refunds/${refund.uid}`);
    await expect(
      page.getByRole("heading", { level: 1, name: refund.refund_id }),
    ).toBeVisible();
    await expect(
      page.getByText("Customer request", { exact: true }),
    ).toBeVisible();
    await expect(page.getByText(/KRW\s15,000/)).toBeVisible();

    await page.getByRole("link", { name: "Refunds for this payment" }).click();

    await expect(page).toHaveURL(new RegExp(`payment_id=${refund.payment_id}`));
    await expect(pageSummary(page, "1–1 of 1")).toBeVisible();
  });

  test("unknown refund shows the not found state", async ({ page }) => {
    await page.goto("/refunds/99999999-9999-4999-8999-999999999999");

    await expect(page.getByText("Refund not found")).toBeVisible();
  });
});
