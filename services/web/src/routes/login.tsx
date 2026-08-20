import { createFileRoute } from "@tanstack/react-router";
import { LoginForm } from "@/domains/auth";

export const Route = createFileRoute("/login")({
  component: LoginPage,
});

function LoginPage() {
  return <LoginForm />;
}
