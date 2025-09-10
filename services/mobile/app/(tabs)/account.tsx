import { Button } from "@bookstore/design/native/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@bookstore/design/native/card";
import { Input } from "@bookstore/design/native/input";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "expo-router";
import { useState } from "react";
import {
  ActivityIndicator,
  Pressable,
  RefreshControl,
  ScrollView,
  Text,
  View,
} from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import {
  customerQueryOptions,
  pointBalanceQueryOptions,
  pointHistoryQueryOptions,
  useUpdateCustomer,
} from "@/domains/account";
import { useAuthStore, useCurrentUserUid } from "@/domains/auth";
import { ErrorState } from "@/shared/components/ErrorState";
import { useToast } from "@/shared/toast";

const TABS = ["profile", "points", "history"] as const;
type Tab = (typeof TABS)[number];

export default function AccountScreen() {
  const [tab, setTab] = useState<Tab>("profile");
  const insets = useSafeAreaInsets();
  const queryClient = useQueryClient();
  const router = useRouter();
  const clearSession = useAuthStore((s) => s.clear);
  const userUid = useCurrentUserUid();
  const [refreshing, setRefreshing] = useState(false);

  const onRefresh = async () => {
    setRefreshing(true);
    try {
      await Promise.all([
        queryClient.invalidateQueries({
          queryKey: ["customer", userUid],
        }),
        queryClient.invalidateQueries({
          queryKey: ["points", userUid],
        }),
      ]);
    } finally {
      setRefreshing(false);
    }
  };

  const handleSignOut = () => {
    clearSession();
    queryClient.clear();
    router.replace("/login");
  };

  return (
    <ScrollView
      className="flex-1 bg-background"
      contentContainerClassName="gap-4 p-5"
      contentContainerStyle={{ paddingBottom: insets.bottom + 20 }}
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
      }
    >
      <Text className="text-2xl font-bold text-foreground">My account</Text>
      <View className="flex-row rounded-md bg-muted p-1">
        {TABS.map((t) => (
          <Pressable
            key={t}
            onPress={() => setTab(t)}
            className={
              tab === t
                ? "flex-1 items-center rounded-md bg-background px-3 py-1.5"
                : "flex-1 items-center px-3 py-1.5"
            }
          >
            <Text
              className={
                tab === t
                  ? "text-sm font-medium capitalize text-foreground"
                  : "text-sm font-medium capitalize text-muted-foreground"
              }
            >
              {t}
            </Text>
          </Pressable>
        ))}
      </View>

      {tab === "profile" ? <ProfileSection /> : null}
      {tab === "points" ? <PointsSection /> : null}
      {tab === "history" ? <HistorySection /> : null}

      <Button variant="outline" onPress={handleSignOut}>
        Sign out
      </Button>
    </ScrollView>
  );
}

function ProfileSection() {
  const userUid = useCurrentUserUid();
  const {
    data: customer,
    isLoading,
    isError,
    refetch,
  } = useQuery(customerQueryOptions(userUid));
  const update = useUpdateCustomer(userUid);

  if (isError) {
    return (
      <Card>
        <ErrorState
          className="items-center py-8"
          message="Couldn't load your profile."
          onRetry={() => refetch()}
        />
      </Card>
    );
  }

  return (
    <ProfileForm
      key={customer?.uid ?? "loading"}
      loading={isLoading}
      initial={customer}
      mutation={update}
    />
  );
}

function ProfileForm({
  loading,
  initial,
  mutation,
}: {
  loading: boolean;
  initial?: {
    first_name?: string;
    last_name?: string;
    email?: string;
    phone?: string;
  };
  mutation: ReturnType<typeof useUpdateCustomer>;
}) {
  const toast = useToast();
  const [firstName, setFirstName] = useState(initial?.first_name ?? "");
  const [lastName, setLastName] = useState(initial?.last_name ?? "");
  const [email, setEmail] = useState(initial?.email ?? "");
  const [phone, setPhone] = useState(initial?.phone ?? "");

  if (loading) {
    return (
      <View className="py-12 items-center">
        <ActivityIndicator />
      </View>
    );
  }

  const handleSave = () => {
    mutation.mutate(
      {
        first_name: firstName,
        last_name: lastName,
        email,
        phone,
      },
      {
        onSuccess: () => toast.success("Profile saved"),
        onError: () => toast.error("Couldn't save profile"),
      },
    );
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Profile</CardTitle>
      </CardHeader>
      <CardContent className="gap-3">
        <Field label="First name" value={firstName} onChange={setFirstName} />
        <Field label="Last name" value={lastName} onChange={setLastName} />
        <Field
          label="Email"
          value={email}
          onChange={setEmail}
          keyboardType="email-address"
          autoCapitalize="none"
        />
        <Field
          label="Phone"
          value={phone}
          onChange={setPhone}
          keyboardType="phone-pad"
        />
        <View className="mt-2 flex-row items-center justify-between">
          <Text className="text-xs text-muted-foreground">
            {mutation.isSuccess ? "Saved." : ""}
          </Text>
          <Button disabled={mutation.isPending} onPress={handleSave}>
            {mutation.isPending ? "Saving..." : "Save changes"}
          </Button>
        </View>
      </CardContent>
    </Card>
  );
}

type FieldProps = {
  label: string;
  value: string;
  onChange: (v: string) => void;
  keyboardType?: "default" | "email-address" | "phone-pad";
  autoCapitalize?: "none" | "sentences";
};

function Field({
  label,
  value,
  onChange,
  keyboardType,
  autoCapitalize,
}: FieldProps) {
  return (
    <View className="gap-1.5">
      <Text className="text-xs font-medium text-muted-foreground">{label}</Text>
      <Input
        value={value}
        onChangeText={onChange}
        keyboardType={keyboardType}
        autoCapitalize={autoCapitalize}
      />
    </View>
  );
}

function PointsSection() {
  const userUid = useCurrentUserUid();
  const { data, isLoading, isError, refetch } = useQuery(
    pointBalanceQueryOptions(userUid),
  );
  return (
    <Card>
      <CardHeader>
        <CardTitle>Reward points</CardTitle>
      </CardHeader>
      <CardContent>
        {isError ? (
          <ErrorState
            className="items-center py-6"
            message="Couldn't load your points."
            onRetry={() => refetch()}
          />
        ) : isLoading ? (
          <ActivityIndicator />
        ) : (
          <View className="flex-row items-center gap-3 rounded-md bg-muted/40 p-4">
            <View className="h-12 w-12 items-center justify-center rounded-full bg-primary">
              <Text className="text-lg font-bold text-primary-foreground">
                ★
              </Text>
            </View>
            <View>
              <Text className="text-xs text-muted-foreground">
                Current balance
              </Text>
              <Text className="text-2xl font-bold text-foreground">
                {(data?.balance ?? 0).toLocaleString()} pt
              </Text>
            </View>
          </View>
        )}
        <Text className="mt-3 text-xs text-muted-foreground">
          Earn 1 point for every $1 spent. Points apply automatically at
          checkout.
        </Text>
      </CardContent>
    </Card>
  );
}

function HistorySection() {
  const userUid = useCurrentUserUid();
  const { data, isLoading, isError, refetch } = useQuery(
    pointHistoryQueryOptions(userUid, { size: 20, page: 0 }),
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>Points history</CardTitle>
      </CardHeader>
      <CardContent>
        {isError ? (
          <ErrorState
            className="items-center py-6"
            message="Couldn't load your history."
            onRetry={() => refetch()}
          />
        ) : isLoading ? (
          <ActivityIndicator />
        ) : !data?.items?.length ? (
          <Text className="py-6 text-center text-sm text-muted-foreground">
            No history yet.
          </Text>
        ) : (
          <View>
            {data.items.map((item) => (
              <View
                key={item.uid}
                className="flex-row items-start justify-between border-t border-border py-3 first:border-t-0"
              >
                <View className="flex-1">
                  <Text className="text-sm font-medium text-foreground">
                    {item.reason ?? item.change_type}
                  </Text>
                  <Text className="text-xs text-muted-foreground">
                    {new Date(item.created_at).toLocaleString()}
                  </Text>
                </View>
                <Text
                  className={
                    item.amount >= 0
                      ? "text-sm font-semibold text-success"
                      : "text-sm font-semibold text-destructive"
                  }
                >
                  {item.amount >= 0 ? "+" : ""}
                  {item.amount.toLocaleString()} pt
                </Text>
              </View>
            ))}
          </View>
        )}
      </CardContent>
    </Card>
  );
}
