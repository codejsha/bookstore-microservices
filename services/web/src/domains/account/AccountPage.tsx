import { Button } from "@bookstore/design/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@bookstore/design/ui/card";
import { Input } from "@bookstore/design/ui/input";
import { Label } from "@bookstore/design/ui/label";
import { Separator } from "@bookstore/design/ui/separator";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@bookstore/design/ui/tabs";
import { useForm } from "@tanstack/react-form";
import { useQuery } from "@tanstack/react-query";
import { Coins, History, Loader2, User } from "lucide-react";
import { useTranslation } from "react-i18next";
import { ErrorState } from "@/shared/components/ErrorState";
import { useToast } from "@/shared/components/toast";
import { useFormat } from "@/shared/hooks/useFormat";
import { useCurrentUserUid } from "@/shared/lib/currentUser";
import { HistorySkeleton, PointsSkeleton } from "./AccountSkeleton";
import { useUpdateCustomer } from "./account-mutations";
import {
  customerQueryOptions,
  pointBalanceQueryOptions,
  pointHistoryQueryOptions,
} from "./account-queries";
import type { PointHistoryItem } from "./types";

const HISTORY_PAGE_SIZE = 20;

const CHANGE_TYPE_KEYS = {
  EARN: "history.changeType.earn",
  SPEND: "history.changeType.spend",
  ADJUST: "history.changeType.adjust",
  EXPIRE: "history.changeType.expire",
} as const satisfies Record<PointHistoryItem["change_type"], string>;

export function AccountPage() {
  const { t } = useTranslation("account");

  return (
    <div className="container mx-auto max-w-3xl space-y-6 px-4 py-8">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{t("title")}</h1>
        <p className="text-sm text-muted-foreground">{t("subtitle")}</p>
      </div>

      <Tabs defaultValue="profile">
        <TabsList>
          <TabsTrigger value="profile">
            <User className="size-4" /> {t("tabs.profile")}
          </TabsTrigger>
          <TabsTrigger value="points">
            <Coins className="size-4" /> {t("tabs.points")}
          </TabsTrigger>
          <TabsTrigger value="history">
            <History className="size-4" /> {t("tabs.history")}
          </TabsTrigger>
        </TabsList>

        <TabsContent value="profile" className="pt-4">
          <ProfileSection />
        </TabsContent>
        <TabsContent value="points" className="pt-4">
          <PointsSection />
        </TabsContent>
        <TabsContent value="history" className="pt-4">
          <HistorySection />
        </TabsContent>
      </Tabs>
    </div>
  );
}

function ProfileSection() {
  const { t } = useTranslation("account");
  const uid = useCurrentUserUid();
  const {
    data: customer,
    isLoading,
    isError,
    refetch,
  } = useQuery(customerQueryOptions(uid));
  const update = useUpdateCustomer(uid);
  const toast = useToast();

  const form = useForm({
    defaultValues: {
      first_name: customer?.first_name ?? "",
      last_name: customer?.last_name ?? "",
      phone: customer?.phone ?? "",
    },
    onSubmit: async ({ value }) => {
      try {
        await update.mutateAsync(value);
        toast.success(t("toast.profileSaved"));
      } catch {
        toast.error(t("toast.profileSaveFailed"));
      }
    },
  });

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <Loader2 className="size-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError) {
    return (
      <ErrorState
        title={t("profile.errorTitle")}
        message={t("profile.errorMessage")}
        onRetry={() => refetch()}
        showHome={false}
      />
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("profile.title")}</CardTitle>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
          className="space-y-4"
        >
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <form.Field name="first_name">
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor={field.name}>{t("profile.firstName")}</Label>
                  <Input
                    id={field.name}
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                </div>
              )}
            </form.Field>
            <form.Field name="last_name">
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor={field.name}>{t("profile.lastName")}</Label>
                  <Input
                    id={field.name}
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                </div>
              )}
            </form.Field>
          </div>

          <div className="space-y-2">
            <Label htmlFor="email">{t("profile.email")}</Label>
            <Input
              id="email"
              type="email"
              value={customer?.email ?? ""}
              readOnly
            />
          </div>

          <form.Field name="phone">
            {(field) => (
              <div className="space-y-2">
                <Label htmlFor={field.name}>{t("profile.phone")}</Label>
                <Input
                  id={field.name}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                />
              </div>
            )}
          </form.Field>

          <Separator />

          <div className="flex items-center justify-end">
            <form.Subscribe selector={(state) => state.isSubmitting}>
              {(isSubmitting) => (
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? t("profile.saving") : t("profile.save")}
                </Button>
              )}
            </form.Subscribe>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}

function PointsSection() {
  const { t } = useTranslation("account");
  const format = useFormat();
  const uid = useCurrentUserUid();
  const { data, isLoading, isError, refetch } = useQuery(
    pointBalanceQueryOptions(uid),
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("points.title")}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <PointsSkeleton />
        ) : isError ? (
          <ErrorState
            title={t("points.errorTitle")}
            message={t("points.errorMessage")}
            onRetry={() => refetch()}
            showHome={false}
          />
        ) : (
          <div className="flex items-center gap-4 rounded-lg border bg-muted/30 px-4 py-6">
            <div className="flex size-12 items-center justify-center rounded-full bg-primary text-primary-foreground">
              <Coins className="size-6" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">
                {t("points.currentBalance")}
              </p>
              <p className="text-3xl font-bold">
                {t("points.amount", {
                  amount: format.number(data?.balance ?? 0),
                })}
              </p>
            </div>
          </div>
        )}
        <p className="mt-4 text-sm text-muted-foreground">
          {t("points.earnInfo")}
        </p>
      </CardContent>
    </Card>
  );
}

function HistorySection() {
  const { t } = useTranslation("account");
  const format = useFormat();
  const uid = useCurrentUserUid();
  const { data, isLoading, isError, refetch } = useQuery(
    pointHistoryQueryOptions(uid, {
      size: HISTORY_PAGE_SIZE,
      page: 0,
    }),
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("history.title")}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <HistorySkeleton />
        ) : isError ? (
          <ErrorState
            title={t("history.errorTitle")}
            message={t("history.errorMessage")}
            onRetry={() => refetch()}
            showHome={false}
          />
        ) : !data?.items?.length ? (
          <p className="py-8 text-center text-sm text-muted-foreground">
            {t("history.empty")}
          </p>
        ) : (
          <ul className="divide-y">
            {data.items.map((item) => (
              <li
                key={item.uid}
                className="flex items-start justify-between gap-4 py-3"
              >
                <div className="min-w-0">
                  <p className="text-sm font-medium">
                    {item.reason ??
                      (CHANGE_TYPE_KEYS[item.change_type]
                        ? t(CHANGE_TYPE_KEYS[item.change_type])
                        : item.change_type)}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {format.dateTime(item.created_at)}
                  </p>
                </div>
                <span
                  className={
                    item.amount >= 0
                      ? "text-sm font-semibold text-green-600"
                      : "text-sm font-semibold text-destructive"
                  }
                >
                  {t("points.amount", {
                    amount: format.number(item.amount, {
                      signDisplay: "always",
                    }),
                  })}
                </span>
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}
