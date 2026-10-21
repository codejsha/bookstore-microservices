import type { Locator, Page } from "@playwright/test";

export function sidebar(page: Page): Locator {
  return page.getByRole("navigation");
}

export function dataRows(scope: Page | Locator): Locator {
  return scope.getByRole("row").filter({ has: scope.getByRole("cell") });
}

export function rowWithCell(scope: Page | Locator, text: string): Locator {
  return scope
    .getByRole("row")
    .filter({ has: scope.getByRole("cell", { name: text, exact: true }) });
}

export function pageSummary(page: Page, text: string): Locator {
  return page.getByText(text, { exact: true });
}

export async function selectFilter(
  page: Page,
  label: string,
  option: string,
): Promise<void> {
  await page.getByRole("combobox", { name: label }).click();
  await page.getByRole("option", { name: option, exact: true }).click();
}
