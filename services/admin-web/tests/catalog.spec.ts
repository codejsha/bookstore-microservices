import { expect, test } from "./fixtures";
import { dataRows, pageSummary, rowWithCell } from "./support/locators";
import { serverError } from "./support/requests";

const WORKS_URL = /\/api\/v1\/admin\/catalog\/works(\?|$)/;
const WORK_URL = /\/api\/v1\/admin\/catalog\/works\/[^/?]+$/;

test.describe("works listing", () => {
  test("lists the first page of works with pagination", async ({ page }) => {
    await page.goto("/catalog");

    await expect(
      page.getByRole("heading", { level: 1, name: "Works" }),
    ).toBeVisible();
    await expect(dataRows(page)).toHaveCount(20);
    await expect(rowWithCell(page, "Work 1")).toContainText("Frank Herbert");
    await expect(rowWithCell(page, "Work 1")).toContainText("Science Fiction");
    await expect(rowWithCell(page, "Work 1")).toContainText("1960");
    await expect(pageSummary(page, "1–20 of 43")).toBeVisible();
    await expect(page.getByText("1 / 3")).toBeVisible();
    await expect(page.getByRole("button", { name: "Previous" })).toBeDisabled();
  });

  test("pagination moves through the pages and updates the url", async ({
    page,
  }) => {
    await page.goto("/catalog");
    const previous = page.getByRole("button", { name: "Previous" });
    const next = page.getByRole("button", { name: "Next" });

    await next.click();

    await expect(page).toHaveURL(/page=1/);
    await expect(pageSummary(page, "21–40 of 43")).toBeVisible();
    await expect(rowWithCell(page, "Work 21")).toBeVisible();
    await expect(previous).toBeEnabled();

    await next.click();

    await expect(pageSummary(page, "41–43 of 43")).toBeVisible();
    await expect(dataRows(page)).toHaveCount(3);
    await expect(next).toBeDisabled();

    await previous.click();

    await expect(pageSummary(page, "21–40 of 43")).toBeVisible();
  });

  test("title filter narrows the list and resets the page", async ({
    page,
  }) => {
    await page.goto("/catalog?page=1");
    await expect(pageSummary(page, "21–40 of 43")).toBeVisible();

    await page.getByRole("textbox", { name: "Filter by title" }).fill("Work 4");

    await expect(page).toHaveURL(/title=Work(\+|%20)4/);
    await expect(page).toHaveURL(/page=0/);
    await expect(dataRows(page)).toHaveCount(5);
    await expect(rowWithCell(page, "Work 4")).toBeVisible();
    await expect(rowWithCell(page, "Work 43")).toBeVisible();
  });

  test("filter without matches shows the empty message", async ({ page }) => {
    await page.goto("/catalog");

    await page
      .getByRole("textbox", { name: "Filter by title" })
      .fill("no such work");

    await expect(page.getByText("No works match this filter.")).toBeVisible();
    await expect(pageSummary(page, "0–0 of 0")).toBeVisible();
  });

  test("title header toggles ascending and descending sort", async ({
    page,
  }) => {
    await page.goto("/catalog");
    const titleHeader = page.getByRole("button", { name: "Title" });

    await titleHeader.click();

    await expect(page).toHaveURL(/sort=title(,|%2C)asc/);
    await expect(dataRows(page).first()).toContainText("Work 1");
    await expect(dataRows(page).nth(1)).toContainText("Work 10");

    await titleHeader.click();

    await expect(page).toHaveURL(/sort=title(,|%2C)desc/);
    await expect(dataRows(page).first()).toContainText("Work 9");
  });

  test("row click opens the work editor", async ({ page, mockApi }) => {
    const work = mockApi.data.works[2];
    await page.goto("/catalog");

    await page.getByRole("cell", { name: work.title, exact: true }).click();

    await expect(page).toHaveURL(`/catalog/${work.uid}`);
    await expect(
      page.getByRole("heading", { level: 1, name: work.title }),
    ).toBeVisible();
    await expect(page.getByText(work.uid, { exact: true })).toBeVisible();
  });

  test("new work button opens the create form", async ({ page }) => {
    await page.goto("/catalog");

    await page.getByRole("button", { name: "New work" }).click();

    await expect(page).toHaveURL("/catalog/new");
    await expect(
      page.getByRole("heading", { level: 1, name: "New work" }),
    ).toBeVisible();
  });
});

test.describe("work creation", () => {
  test("submitting an empty form reports the required fields", async ({
    page,
  }) => {
    await page.goto("/catalog/new");

    await page.getByRole("button", { name: "Create work" }).click();

    await expect(
      page.getByRole("alert").filter({ hasText: "Title is required" }),
    ).toBeVisible();
    await expect(
      page.getByRole("alert").filter({ hasText: "Pick at least one author" }),
    ).toBeVisible();
  });

  test("valid form posts the work and returns to the list", async ({
    page,
    mockApi,
  }) => {
    const author = mockApi.data.authors[1];
    await page.goto("/catalog/new");
    await page
      .getByRole("textbox", { name: "Title", exact: true })
      .fill("Dune");
    await page
      .getByRole("listbox", { name: "Authors" })
      .selectOption({ label: author.name });
    await page
      .getByRole("textbox", { name: "Subjects" })
      .fill("Fantasy, Space opera");
    await page.getByRole("textbox", { name: "First published" }).fill("1965");
    const createRequest = page.waitForRequest(
      (request) => request.method() === "POST" && WORKS_URL.test(request.url()),
    );

    await page.getByRole("button", { name: "Create work" }).click();

    expect((await createRequest).postDataJSON()).toEqual({
      title: "Dune",
      first_publish_date: "1965",
      author_uids: [author.uid],
      subject_names: ["Fantasy", "Space opera"],
    });
    await expect(page).toHaveURL(/\/catalog\?page=0$/);
    await expect(
      page.getByRole("heading", { level: 1, name: "Works" }),
    ).toBeVisible();
  });

  test("failed create keeps the form and shows the error", async ({
    page,
    mockApi,
  }) => {
    await page.route(WORKS_URL, (route) =>
      route.request().method() === "POST"
        ? route.fulfill(serverError())
        : route.fallback(),
    );
    await page.goto("/catalog/new");
    await page
      .getByRole("textbox", { name: "Title", exact: true })
      .fill("Dune");
    await page
      .getByRole("listbox", { name: "Authors" })
      .selectOption({ label: mockApi.data.authors[0].name });

    await page.getByRole("button", { name: "Create work" }).click();

    await expect(page.getByText("Could not create this work")).toBeVisible();
    await expect(page.getByText("Internal Server Error")).toBeVisible();
    await expect(page).toHaveURL("/catalog/new");
    await expect(
      page.getByRole("textbox", { name: "Title", exact: true }),
    ).toHaveValue("Dune");
  });

  test("cancel returns to the list without posting", async ({ page }) => {
    const posts: string[] = [];
    page.on("request", (request) => {
      if (request.method() === "POST") posts.push(request.url());
    });
    await page.goto("/catalog/new");

    await page.getByRole("button", { name: "Cancel" }).click();

    await expect(page).toHaveURL(/\/catalog\?page=0$/);
    expect(posts).toHaveLength(0);
  });
});

test.describe("work editor", () => {
  test("editor is prefilled from the work", async ({ page, mockApi }) => {
    const work = mockApi.data.works[0];
    await page.goto(`/catalog/${work.uid}`);

    await expect(
      page.getByRole("textbox", { name: "Title", exact: true }),
    ).toHaveValue("Work 1");
    await expect(page.getByRole("listbox", { name: "Authors" })).toHaveValues([
      work.authors[0].uid,
    ]);
    await expect(page.getByRole("textbox", { name: "Subjects" })).toHaveValue(
      "Science Fiction",
    );
    await expect(
      page.getByRole("textbox", { name: "First published" }),
    ).toHaveValue("1960");
    await expect(
      page.getByRole("textbox", { name: "Open Library key" }),
    ).toHaveValue("OL0W");
    await expect(
      page.getByRole("textbox", { name: "Description" }),
    ).toHaveValue("A description for work 1.");
    await expect(
      page.getByRole("button", { name: "Save changes" }),
    ).toBeEnabled();
  });

  test("saving sends the full payload and updates the work", async ({
    page,
    mockApi,
  }) => {
    const work = mockApi.data.works[0];
    await page.goto(`/catalog/${work.uid}`);
    await page
      .getByRole("textbox", { name: "Title", exact: true })
      .fill("Dune");
    const updateRequest = page.waitForRequest(
      (request) =>
        request.method() === "PUT" && request.url().endsWith(work.uid),
    );

    await page.getByRole("button", { name: "Save changes" }).click();

    expect((await updateRequest).postDataJSON()).toEqual({
      title: "Dune",
      description: "A description for work 1.",
      first_publish_date: "1960",
      ol_key: "OL0W",
      author_uids: [work.authors[0].uid],
      subject_names: ["Science Fiction"],
    });
    await expect(
      page.getByRole("heading", { level: 1, name: "Dune" }),
    ).toBeVisible();

    await page.getByRole("button", { name: "Cancel" }).click();

    await expect(page).toHaveURL(/\/catalog\?page=0$/);
    await expect(rowWithCell(page, "Dune")).toBeVisible();
    await expect(rowWithCell(page, "Work 1")).toHaveCount(0);
  });

  test("failed save shows the error and can be dismissed", async ({
    page,
    mockApi,
  }) => {
    const work = mockApi.data.works[0];
    await page.route(WORK_URL, (route) =>
      route.request().method() === "PUT"
        ? route.fulfill(serverError())
        : route.fallback(),
    );
    await page.goto(`/catalog/${work.uid}`);

    await page.getByRole("button", { name: "Save changes" }).click();

    await expect(page.getByText("Could not save this work")).toBeVisible();
    await expect(page.getByText("Internal Server Error")).toBeVisible();

    await page.getByRole("button", { name: "Try again" }).click();

    await expect(page.getByText("Could not save this work")).toHaveCount(0);
    await expect(
      page.getByRole("heading", { level: 1, name: "Work 1" }),
    ).toBeVisible();
  });

  test("unknown work shows the not found page", async ({ page }) => {
    await page.goto("/catalog/22222222-2222-2222-2222-999999999999");

    await expect(page.getByText("Page not found")).toBeVisible({
      timeout: 15_000,
    });
  });
});
