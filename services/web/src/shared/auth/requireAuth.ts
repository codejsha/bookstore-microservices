import { redirect } from "@tanstack/react-router";
import { userManager } from "./userManager";

export async function requireAuth(): Promise<void> {
  const user = await userManager.getUser();
  if (!user || user.expired) {
    throw redirect({ to: "/login" });
  }
}
