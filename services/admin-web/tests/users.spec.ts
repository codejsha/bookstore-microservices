import { expect, test } from "./fixtures";
import {
  dataRows,
  pageSummary,
  rowWithCell,
  selectFilter,
} from "./support/locators";

test.describe("users listing", () => {
  test("lists the first page of users", async ({ page }) => {
    await page.goto("/users");
    const first = rowWithCell(page, "ada.lovelace0@example.com");

    await expect(
      page.getByRole("heading", { level: 1, name: "Users" }),
    ).toBeVisible();
    await expect(dataRows(page)).toHaveCount(20);
    await expect(pageSummary(page, "1–20 of 32")).toBeVisible();
    await expect(first).toContainText("Ada");
    await expect(first).toContainText("Lovelace");
    await expect(first).toContainText("ACTIVE");
    await expect(first).toContainText("USER, STAFF");
  });

  test("email filter narrows the list", async ({ page }) => {
    await page.goto("/users");

    await page.getByRole("textbox", { name: "Email" }).fill("hopper");

    await expect(page).toHaveURL(/email=hopper/);
    await expect(pageSummary(page, "1–7 of 7")).toBeVisible();
    await expect(rowWithCell(page, "grace.hopper1@example.com")).toBeVisible();
  });

  test("name filter matches first or last names", async ({ page }) => {
    await page.goto("/users");

    await page.getByRole("textbox", { name: "Name" }).fill("Grace");

    await expect(pageSummary(page, "1–6 of 6")).toBeVisible();
    await expect(dataRows(page).first()).toContainText("Grace");
  });

  test("status filter narrows the list and can be cleared", async ({
    page,
  }) => {
    await page.goto("/users");

    await selectFilter(page, "Status", "SUSPENDED");

    await expect(page).toHaveURL(/status=SUSPENDED/);
    await expect(page.getByRole("combobox", { name: "Status" })).toContainText(
      "SUSPENDED",
    );
    await expect(pageSummary(page, "1–3 of 3")).toBeVisible();
    await expect(dataRows(page).first()).toContainText("SUSPENDED");

    await selectFilter(page, "Status", "All statuses");

    await expect(page).not.toHaveURL(/status=/);
    await expect(pageSummary(page, "1–20 of 32")).toBeVisible();
  });

  test("filter without matches shows the empty message", async ({ page }) => {
    await page.goto("/users");

    await page.getByRole("textbox", { name: "Email" }).fill("nobody");

    await expect(page.getByText("No users match these filters.")).toBeVisible();
  });

  test("email header sorts ascending then descending", async ({ page }) => {
    await page.goto("/users");
    const emailHeader = page.getByRole("button", { name: "Email" });

    await emailHeader.click();

    await expect(page).toHaveURL(/sort=email(,|%2C)asc/);
    await expect(dataRows(page).first()).toContainText(
      "ada.dijkstra24@example.com",
    );

    await emailHeader.click();

    await expect(page).toHaveURL(/sort=email(,|%2C)desc/);
    await expect(dataRows(page).first()).toContainText(
      "margaret.turing17@example.com",
    );
  });

  test("row click opens the user detail", async ({ page, mockApi }) => {
    const user = mockApi.data.users[0];
    await page.goto("/users");

    await page.getByRole("cell", { name: user.email, exact: true }).click();

    await expect(page).toHaveURL(`/users/${user.uid}`);
    await expect(
      page.getByRole("heading", { level: 1, name: "Ada Lovelace" }),
    ).toBeVisible();
    await expect(page.getByText(user.uid, { exact: true })).toBeVisible();
    await expect(page.getByText("USER, STAFF")).toBeVisible();
  });
});

test.describe("user detail", () => {
  test("view orders link opens the orders filtered by the user", async ({
    page,
    mockApi,
  }) => {
    const user = mockApi.data.users[0];
    await page.goto(`/users/${user.uid}`);

    await page.getByRole("link", { name: "View orders" }).click();

    await expect(page).toHaveURL(/\/orders\?/);
    await expect(page).toHaveURL(new RegExp(`user_uid=${user.uid}`));
    await expect(page.getByRole("textbox", { name: "User uid" })).toHaveValue(
      user.uid,
    );
    await expect(pageSummary(page, "1–4 of 4")).toBeVisible();
  });

  test("back link returns to the listing", async ({ page, mockApi }) => {
    await page.goto(`/users/${mockApi.data.users[0].uid}`);

    await page.getByRole("link", { name: "Back to users" }).click();

    await expect(page).toHaveURL(/\/users\?page=0$/);
  });

  test("unknown user shows the not found state", async ({ page }) => {
    await page.goto("/users/55555555-5555-4555-8555-999999999999");

    await expect(page.getByText("User not found")).toBeVisible();
    await expect(page.getByText("No user exists with this id.")).toBeVisible();
  });
});
