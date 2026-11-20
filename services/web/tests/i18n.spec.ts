import {
  DEFAULT_LOCALE,
  LOCALES,
  type Locale,
} from "../src/shared/i18n/locales";
import { expect, test } from "./fixtures";
import { localePath } from "./support/paths";

const COPY: Record<Locale, { books: string; menu: string }> = {
  en: { books: "Books", menu: "Language" },
  ko: { books: "도서", menu: "언어" },
  ja: { books: "書籍", menu: "言語" },
};

const TRANSLATED_LOCALES = LOCALES.filter(
  (locale) => locale.code !== DEFAULT_LOCALE,
);

test.describe("locale routing", () => {
  test("unprefixed home redirects to the default locale", async ({ page }) => {
    await page.goto("/");

    await expect(page).toHaveURL(localePath("/"));
    await expect(page.locator("html")).toHaveAttribute("lang", "en");
    await expect(
      page.getByRole("heading", {
        level: 1,
        name: "Discover Your Next Great Read",
      }),
    ).toBeVisible();
  });

  test("unprefixed deep link keeps its path, search and hash", async ({
    page,
  }) => {
    await page.goto("/books/work-1?ref=mail#reviews");

    await expect(page).toHaveURL(
      `${localePath("/books/work-1")}?ref=mail#reviews`,
    );
    await expect(
      page.getByRole("heading", { level: 1, name: "Clean Architecture" }),
    ).toBeVisible();
  });

  test("unknown first segment lands on not found under the default locale", async ({
    page,
  }) => {
    await page.goto("/fr/books");

    await expect(page).toHaveURL(localePath("/fr/books"));
    await expect(page.getByText("Page not found")).toBeVisible();
    await expect(page.getByRole("banner")).toHaveCount(1);
  });

  for (const segment of ["kr", "jp"]) {
    test(`country code ${segment} is not a locale prefix`, async ({ page }) => {
      await page.goto(`/${segment}/books`);

      await expect(page).toHaveURL(localePath(`/${segment}/books`));
      await expect(page.locator("html")).toHaveAttribute("lang", "en");
      await expect(page.getByText("Page not found")).toBeVisible();
    });
  }

  test("callback stays unprefixed", async ({ page }) => {
    await page.goto("/callback");

    await expect(page).toHaveURL("/callback");
    await expect(page.getByText("Completing sign-in…")).toBeVisible();
  });

  for (const { code, name } of TRANSLATED_LOCALES) {
    const copy = COPY[code];

    test(`${code} prefix renders ${name}`, async ({ page }) => {
      await page.goto(localePath("/books", code));

      await expect(page.locator("html")).toHaveAttribute("lang", code);
      await expect(
        page.getByRole("heading", { level: 1, name: copy.books }),
      ).toBeVisible();
      await expect(page.getByRole("button", { name: copy.menu })).toBeVisible();
    });

    test(`switcher moves between en and ${code} on the same page`, async ({
      page,
    }) => {
      await page.goto(`${localePath("/books/work-1")}?ref=mail#reviews`);

      await page.getByRole("button", { name: COPY.en.menu }).click();
      await page.getByRole("menuitemradio", { name }).click();

      await expect(page).toHaveURL(
        `${localePath("/books/work-1", code)}?ref=mail#reviews`,
      );
      await expect(page.locator("html")).toHaveAttribute("lang", code);
      await expect(
        page.getByRole("heading", { level: 1, name: "Clean Architecture" }),
      ).toBeVisible();
      await expect(page.getByRole("menu")).toBeHidden();

      await page.getByRole("button", { name: copy.menu }).click();
      await page.getByRole("menuitemradio", { name: "English" }).click();

      await expect(page).toHaveURL(
        `${localePath("/books/work-1")}?ref=mail#reviews`,
      );
      await expect(page.locator("html")).toHaveAttribute("lang", "en");
    });

    test(`in-app navigation keeps the ${code} prefix`, async ({ page }) => {
      await page.goto(localePath("/", code));

      await page
        .getByRole("complementary")
        .first()
        .getByRole("link", { name: copy.books })
        .click();

      await expect(page).toHaveURL(localePath("/books", code));

      await page.getByRole("link", { name: "Clean Architecture" }).click();

      await expect(page).toHaveURL(localePath("/books/work-1", code));
      await expect(page.locator("html")).toHaveAttribute("lang", code);
    });

    test(`${code} visit is not remembered for unprefixed urls`, async ({
      page,
    }) => {
      await page.goto(localePath("/books", code));
      await expect(page.locator("html")).toHaveAttribute("lang", code);

      await page.goto("/books");

      await expect(page).toHaveURL(localePath("/books"));
      await expect(page.locator("html")).toHaveAttribute("lang", "en");
      await expect(
        page.getByRole("heading", { level: 1, name: "Books" }),
      ).toBeVisible();
    });

    test(`protected route under ${code} redirects to its login`, async ({
      page,
    }) => {
      await page.goto(localePath("/cart", code));

      await expect(page).toHaveURL(localePath("/login", code));
    });
  }
});

test.describe("locale routing with a korean browser", () => {
  test.use({ locale: "ko-KR" });

  test("unprefixed home ignores the browser language", async ({ page }) => {
    await page.goto("/");

    await expect(page).toHaveURL(localePath("/"));
    await expect(page.locator("html")).toHaveAttribute("lang", "en");
  });

  test("unprefixed deep link ignores the browser language", async ({
    page,
  }) => {
    await page.goto("/books");

    await expect(page).toHaveURL(localePath("/books"));
    await expect(
      page.getByRole("heading", { level: 1, name: "Books" }),
    ).toBeVisible();
  });

  test("japanese prefix wins over the browser language", async ({ page }) => {
    await page.goto(localePath("/books", "ja"));

    await expect(page.locator("html")).toHaveAttribute("lang", "ja");
    await expect(
      page.getByRole("heading", { level: 1, name: COPY.ja.books }),
    ).toBeVisible();
  });
});
