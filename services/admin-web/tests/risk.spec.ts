import { expect, test } from "./fixtures";
import { dataRows, rowWithCell } from "./support/locators";
import { serverError } from "./support/requests";

const NEW_UID = "aaaaaaaa-0000-4000-8000-000000000001";

test.describe("risk flags", () => {
  test("lists the flagged principals with their level", async ({
    page,
    mockApi,
  }) => {
    const [restricted, blocked] = mockApi.data.riskEntries;
    await page.goto("/risk");

    await expect(dataRows(page)).toHaveCount(2);
    await expect(rowWithCell(page, restricted.user_uid)).toContainText(
      "restrict",
    );
    await expect(rowWithCell(page, restricted.user_uid)).toContainText(
      restricted.reason,
    );
    await expect(rowWithCell(page, blocked.user_uid)).toContainText("block");
  });

  test("flag button stays disabled until a uid and reason are entered", async ({
    page,
  }) => {
    await page.goto("/risk");
    const flag = page.getByRole("button", { name: "Flag", exact: true });

    await expect(flag).toBeDisabled();

    await page.getByRole("textbox", { name: "User uid" }).fill(NEW_UID);

    await expect(flag).toBeDisabled();

    await page.getByRole("textbox", { name: "Reason" }).fill("abuse");

    await expect(flag).toBeEnabled();
  });

  test("flagging a principal sends the request and adds a row", async ({
    page,
  }) => {
    await page.goto("/risk");
    await page.getByRole("textbox", { name: "User uid" }).fill(NEW_UID);
    await page.getByLabel("Level").selectOption("block");
    await page
      .getByRole("textbox", { name: "Reason" })
      .fill("chargeback abuse");
    await page.getByLabel("TTL (days)").fill("2");
    const flagRequest = page.waitForRequest(
      (request) =>
        request.method() === "PUT" &&
        request.url().endsWith(`/api/v1/admin/risk/${NEW_UID}`),
    );

    await page.getByRole("button", { name: "Flag", exact: true }).click();

    expect((await flagRequest).postDataJSON()).toEqual({
      level: "block",
      reason: "chargeback abuse",
      ttl_seconds: 172_800,
    });
    await expect(dataRows(page)).toHaveCount(3);
    await expect(rowWithCell(page, NEW_UID)).toContainText("block");
    await expect(rowWithCell(page, NEW_UID)).toContainText("chargeback abuse");
    await expect(page.getByRole("textbox", { name: "User uid" })).toHaveValue(
      "",
    );
    await expect(page.getByRole("textbox", { name: "Reason" })).toHaveValue("");
  });

  test("unflagging removes the row", async ({ page, mockApi }) => {
    const [restricted, blocked] = mockApi.data.riskEntries;
    await page.goto("/risk");
    const unflagRequest = page.waitForRequest(
      (request) =>
        request.method() === "DELETE" &&
        request.url().endsWith(`/api/v1/admin/risk/${restricted.user_uid}`),
    );

    await rowWithCell(page, restricted.user_uid)
      .getByRole("button", { name: "Unflag" })
      .click();

    await unflagRequest;
    await expect(dataRows(page)).toHaveCount(1);
    await expect(rowWithCell(page, blocked.user_uid)).toBeVisible();
  });

  test("failed flag shows the error message", async ({ page }) => {
    await page.route(/\/api\/v1\/admin\/risk\/[^/]+$/, (route) =>
      route.request().method() === "PUT"
        ? route.fulfill(serverError())
        : route.fallback(),
    );
    await page.goto("/risk");
    await page.getByRole("textbox", { name: "User uid" }).fill(NEW_UID);
    await page.getByRole("textbox", { name: "Reason" }).fill("abuse");

    await page.getByRole("button", { name: "Flag", exact: true }).click();

    await expect(page.getByText("Internal Server Error")).toBeVisible();
    await expect(dataRows(page)).toHaveCount(2);
  });
});
