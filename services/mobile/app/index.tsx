import { Redirect } from "expo-router";
import { isSessionValid, useAuthStore } from "@/domains/auth";

export default function Index() {
  const authed = useAuthStore(isSessionValid);
  return <Redirect href={authed ? "/(tabs)" : "/login"} />;
}
