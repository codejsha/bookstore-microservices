import type { Locator, Page } from "@playwright/test";

export function browseSidebar(page: Page): Locator {
  return page.getByRole("complementary", { name: "Browse" });
}

export function accountSidebar(page: Page): Locator {
  return page.getByRole("complementary", { name: "Account" });
}

export function bookItems(page: Page): Locator {
  return page.getByRole("list", { name: "Books" }).getByRole("listitem");
}

export function cartItems(page: Page): Locator {
  return page.getByRole("list", { name: "Cart items" }).getByRole("listitem");
}

export function cartItem(page: Page, name: string): Locator {
  return cartItems(page).filter({ hasText: name });
}

export function wishlistItems(page: Page): Locator {
  return page
    .getByRole("list", { name: "Wishlist items" })
    .getByRole("listitem");
}

export function toast(page: Page, message: string): Locator {
  return page.getByRole("status").filter({ hasText: message });
}
