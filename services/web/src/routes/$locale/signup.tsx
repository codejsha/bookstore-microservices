import { createFileRoute } from "@tanstack/react-router";
import { SignUpForm } from "@/domains/auth";

export const Route = createFileRoute("/$locale/signup")({
  component: SignUpPage,
});

function SignUpPage() {
  return <SignUpForm />;
}
