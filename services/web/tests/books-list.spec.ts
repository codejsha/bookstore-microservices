import { expect, test } from "./fixtures";
import { bookItems } from "./support/locators";
import { localePath } from "./support/paths";
import { apiQuery, serverError } from "./support/requests";

const WORKS_PATH = "/api/v1/works";
const WORKS_URL = /\/api\/v1\/works\?/;

test.describe("books listing", () => {
  test("shows the first page with filter, search and sort controls", async ({
    page,
  }) => {
    await page.goto(localePath("/books"));

    await expect(
      page.getByRole("heading", { level: 1, name: "Books" }),
    ).toBeVisible();
    await expect(page.getByText("24 books found")).toBeVisible();
    await expect(
      page.getByRole("button", { name: "All", exact: true }),
    ).toHaveAttribute("aria-pressed", "true");
    for (const subject of [
      "Fiction",
      "Non-Fiction",
      "Science",
      "Technology",
      "History",
    ]) {
      await expect(
        page.getByRole("button", { name: subject, exact: true }),
      ).toHaveAttribute("aria-pressed", "false");
    }
    await expect(page.getByRole("combobox", { name: "Sort by" })).toHaveValue(
      "title,asc",
    );
    await expect(
      page.getByRole("button", { name: "Grid view" }),
    ).toHaveAttribute("aria-pressed", "true");
    await expect(bookItems(page)).toHaveCount(12);
    await expect(bookItems(page).first()).toContainText(
      "A Philosophy of Software Design",
    );
  });

  test("search filters by title through a single debounced request", async ({
    page,
  }) => {
    await page.goto(localePath("/books"));
    await expect(bookItems(page)).toHaveCount(12);
    const titleQueries: string[] = [];
    page.on("request", (request) => {
      const title = apiQuery(request, WORKS_PATH)?.get("title");
      if (title !== null && title !== undefined) titleQueries.push(title);
    });

    await page
      .getByRole("textbox", { name: "Search by title" })
      .pressSequentially("clean", { delay: 50 });

    await expect(bookItems(page)).toHaveCount(2);
    await expect(page.getByText("2 books found")).toBeVisible();
    await expect(bookItems(page).first()).toContainText("Clean Architecture");
    await expect(bookItems(page).last()).toContainText("Clean Code");
    expect(titleQueries).toEqual(["clean"]);
  });

  test("search without matches shows the empty state", async ({ page }) => {
    await page.goto(localePath("/books"));
    await expect(bookItems(page)).toHaveCount(12);

    await page.getByRole("textbox", { name: "Search by title" }).fill("zzzz");

    await expect(page.getByText("No books found")).toBeVisible();
    await expect(
      page.getByText("Try adjusting your search or filters."),
    ).toBeVisible();
    await expect(bookItems(page)).toHaveCount(0);
  });

  test("subject filter narrows the list and resets with All", async ({
    page,
  }) => {
    await page.goto(localePath("/books"));
    await expect(bookItems(page)).toHaveCount(12);
    const fiction = page.getByRole("button", { name: "Fiction", exact: true });
    const all = page.getByRole("button", { name: "All", exact: true });
    const filtered = page.waitForRequest(
      (request) =>
        apiQuery(request, WORKS_PATH)?.get("subject_uid") === "sub-1",
    );

    await fiction.click();

    await filtered;
    await expect(fiction).toHaveAttribute("aria-pressed", "true");
    await expect(all).toHaveAttribute("aria-pressed", "false");
    await expect(page.getByText("5 books found")).toBeVisible();
    await expect(bookItems(page)).toHaveCount(5);

    await all.click();

    await expect(all).toHaveAttribute("aria-pressed", "true");
    await expect(bookItems(page)).toHaveCount(12);
  });

  test("sort selection reorders the list", async ({ page }) => {
    await page.goto(localePath("/books"));
    await expect(bookItems(page).first()).toContainText(
      "A Philosophy of Software Design",
    );

    await page
      .getByRole("combobox", { name: "Sort by" })
      .selectOption({ label: "Title Z-A" });

    await expect(bookItems(page).first()).toContainText(
      "Working Effectively with Legacy Code",
    );
  });

  test("view toggle switches between grid and list", async ({ page }) => {
    await page.goto(localePath("/books"));
    const grid = page.getByRole("button", { name: "Grid view" });
    const list = page.getByRole("button", { name: "List view" });
    await expect(bookItems(page)).toHaveCount(12);

    await list.click();

    await expect(list).toHaveAttribute("aria-pressed", "true");
    await expect(grid).toHaveAttribute("aria-pressed", "false");
    await expect(bookItems(page)).toHaveCount(12);
    await expect(bookItems(page).first()).toContainText(
      "A Philosophy of Software Design",
    );

    await grid.click();

    await expect(grid).toHaveAttribute("aria-pressed", "true");
  });

  test("scrolling to the end loads the remaining books", async ({ page }) => {
    await page.goto(localePath("/books"));
    await expect(bookItems(page)).toHaveCount(12);

    await bookItems(page).last().scrollIntoViewIfNeeded();

    await expect(bookItems(page)).toHaveCount(24);
    await expect(page.getByText("All 24 books loaded")).toBeVisible();
  });

  test("book card opens the detail page", async ({ page }) => {
    await page.goto(localePath("/books"));

    await bookItems(page).first().getByRole("link").click();

    await expect(page).toHaveURL(/\/books\/work-\d+$/);
    await expect(page.getByRole("heading", { level: 1 })).toHaveText(
      "A Philosophy of Software Design",
    );
  });

  test("server failure shows an error state that can be retried", async ({
    page,
  }) => {
    await page.route(WORKS_URL, (route) => route.fulfill(serverError()));
    await page.goto(localePath("/books"));

    await expect(page.getByText("Couldn't load books")).toBeVisible({
      timeout: 15_000,
    });

    await page.unroute(WORKS_URL);
    await page.getByRole("button", { name: "Try again" }).click();

    await expect(bookItems(page)).toHaveCount(12);
  });
});
