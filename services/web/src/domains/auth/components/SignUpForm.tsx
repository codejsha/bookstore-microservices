import { Button } from "@bookstore/design/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@bookstore/design/ui/card";
import { Link } from "@tanstack/react-router";
import { useAuth } from "react-oidc-context";

export function SignUpForm() {
  const auth = useAuth();

  const register = () => {
    auth.signinRedirect({ extraQueryParams: { kc_action: "register" } });
  };

  return (
    <div className="flex min-h-[calc(100vh-3.5rem)] flex-col items-center justify-center px-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Create account</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-2">
          <Button onClick={register} disabled={auth.isLoading}>
            Continue with Keycloak
          </Button>
          {auth.error && (
            <p className="text-sm text-destructive">{auth.error.message}</p>
          )}
        </CardContent>
      </Card>
      <div className="mt-4 flex items-center justify-center gap-4 text-sm">
        <Link
          to="/login"
          className="text-muted-foreground hover:text-foreground"
        >
          Already have an account? Login
        </Link>
      </div>
    </div>
  );
}
