import { redirect } from "@tanstack/react-router";
import { userManager } from "./userManager";

export async function requireAuth(locale: string): Promise<void> {
  const user = await userManager.getUser();
  if (!user || user.expired) {
    throw redirect({ to: "/$locale/login", params: { locale } });
  }
}
