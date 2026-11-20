import { Button } from "@bookstore/design/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@bookstore/design/ui/card";
import { Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { useAuth } from "react-oidc-context";
import { signinState } from "@/shared/auth/signinState";
import { useLocale } from "@/shared/i18n/useLocale";

export function SignUpForm() {
  const { t } = useTranslation("auth");
  const auth = useAuth();
  const locale = useLocale();

  const register = () => {
    auth.signinRedirect({
      extraQueryParams: { kc_action: "register" },
      state: signinState(locale),
    });
  };

  return (
    <div className="flex min-h-[calc(100vh-3.5rem)] flex-col items-center justify-center px-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">{t("signUp.title")}</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-2">
          <Button onClick={register} disabled={auth.isLoading}>
            {t("signUp.continue")}
          </Button>
          {auth.error && (
            <p className="text-sm text-destructive">{auth.error.message}</p>
          )}
        </CardContent>
      </Card>
      <div className="mt-4 flex items-center justify-center gap-4 text-sm">
        <Link
          to="/$locale/login"
          params={{ locale }}
          className="text-muted-foreground hover:text-foreground"
        >
          {t("signUp.haveAccount")}
        </Link>
      </div>
    </div>
  );
}
