import { buttonVariants } from "@bookstore/design/ui/button";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/unauthorized")({
  component: UnauthorizedPage,
});

function UnauthorizedPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-3 p-8 text-center">
      <h1 className="font-semibold text-lg">Not authorized</h1>
      <p className="max-w-md text-muted-foreground text-sm">
        Your account does not have the STAFF role. Ask a manager to grant it, or
        sign in with a different account.
      </p>
      <a
        href="/oauth2/sign_out"
        className={buttonVariants({ variant: "outline" })}
      >
        Sign out
      </a>
    </div>
  );
}
