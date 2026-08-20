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
import { ErrorState } from "@/shared/components/ErrorState";
import { useToast } from "@/shared/components/toast";
import { useCurrentUserUid } from "@/shared/lib/currentUser";
import { HistorySkeleton, PointsSkeleton } from "./AccountSkeleton";
import { useUpdateCustomer } from "./account-mutations";
import {
  customerQueryOptions,
  pointBalanceQueryOptions,
  pointHistoryQueryOptions,
} from "./account-queries";

const HISTORY_PAGE_SIZE = 20;

export function AccountPage() {
  return (
    <div className="container mx-auto max-w-3xl space-y-6 px-4 py-8">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">My Account</h1>
        <p className="text-sm text-muted-foreground">
          Manage your profile, rewards, and activity.
        </p>
      </div>

      <Tabs defaultValue="profile">
        <TabsList>
          <TabsTrigger value="profile">
            <User className="size-4" /> Profile
          </TabsTrigger>
          <TabsTrigger value="points">
            <Coins className="size-4" /> Points
          </TabsTrigger>
          <TabsTrigger value="history">
            <History className="size-4" /> History
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
      email: customer?.email ?? "",
      phone: customer?.phone ?? "",
    },
    onSubmit: async ({ value }) => {
      try {
        await update.mutateAsync(value);
        toast.success("Profile saved");
      } catch {
        toast.error("Couldn't save profile");
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
        title="Couldn't load your profile"
        message="There was a problem loading your account details."
        onRetry={() => refetch()}
        showHome={false}
      />
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Profile</CardTitle>
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
                  <Label htmlFor={field.name}>First name</Label>
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
                  <Label htmlFor={field.name}>Last name</Label>
                  <Input
                    id={field.name}
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                </div>
              )}
            </form.Field>
          </div>

          <form.Field
            name="email"
            validators={{
              onBlur: ({ value }) => {
                if (!value) return undefined;
                if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value))
                  return "Invalid email address";
                return undefined;
              },
            }}
          >
            {(field) => (
              <div className="space-y-2">
                <Label htmlFor={field.name}>Email</Label>
                <Input
                  id={field.name}
                  type="email"
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                />
                {field.state.meta.errors.length > 0 && (
                  <p className="text-sm text-destructive">
                    {field.state.meta.errors[0]}
                  </p>
                )}
              </div>
            )}
          </form.Field>

          <form.Field name="phone">
            {(field) => (
              <div className="space-y-2">
                <Label htmlFor={field.name}>Phone</Label>
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
                  {isSubmitting ? "Saving..." : "Save changes"}
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
  const uid = useCurrentUserUid();
  const { data, isLoading, isError, refetch } = useQuery(
    pointBalanceQueryOptions(uid),
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>Reward points</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <PointsSkeleton />
        ) : isError ? (
          <ErrorState
            title="Couldn't load points"
            message="There was a problem loading your reward balance."
            onRetry={() => refetch()}
            showHome={false}
          />
        ) : (
          <div className="flex items-center gap-4 rounded-lg border bg-muted/30 px-4 py-6">
            <div className="flex size-12 items-center justify-center rounded-full bg-primary text-primary-foreground">
              <Coins className="size-6" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Current balance</p>
              <p className="text-3xl font-bold">
                {(data?.balance ?? 0).toLocaleString()} pt
              </p>
            </div>
          </div>
        )}
        <p className="mt-4 text-sm text-muted-foreground">
          Earn 1 point for every $1 spent. Points apply automatically at
          checkout.
        </p>
      </CardContent>
    </Card>
  );
}

function HistorySection() {
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
        <CardTitle>Points history</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <HistorySkeleton />
        ) : isError ? (
          <ErrorState
            title="Couldn't load history"
            message="There was a problem loading your points history."
            onRetry={() => refetch()}
            showHome={false}
          />
        ) : !data?.items?.length ? (
          <p className="py-8 text-center text-sm text-muted-foreground">
            No history yet.
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
                    {item.reason ?? item.change_type}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {new Date(item.created_at).toLocaleString()}
                  </p>
                </div>
                <span
                  className={
                    item.amount >= 0
                      ? "text-sm font-semibold text-green-600"
                      : "text-sm font-semibold text-destructive"
                  }
                >
                  {item.amount >= 0 ? "+" : ""}
                  {item.amount.toLocaleString()} pt
                </span>
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}
