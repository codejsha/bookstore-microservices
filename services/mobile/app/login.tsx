import { Button } from "@bookstore/design/native/button";
import { Card } from "@bookstore/design/native/card";
import { Link, useRouter } from "expo-router";
import { useEffect } from "react";
import { Text, View } from "react-native";
import { isSessionValid, useAuthStore, useKeycloakAuth } from "@/domains/auth";

export default function LoginScreen() {
  const router = useRouter();
  const authed = useAuthStore(isSessionValid);
  const { request, promptAsync } = useKeycloakAuth();

  useEffect(() => {
    if (authed) router.replace("/(tabs)");
  }, [authed, router]);

  return (
    <View className="flex-1 justify-center bg-background px-5">
      <Card className="gap-4">
        <View className="items-center">
          <Text className="text-2xl font-bold text-foreground">
            Welcome back
          </Text>
          <Text className="mt-1 text-sm text-muted-foreground">
            Sign in with Keycloak
          </Text>
        </View>

        <Button disabled={!request} onPress={() => promptAsync()}>
          Continue with Keycloak
        </Button>

        <View className="flex-row justify-center gap-2">
          <Text className="text-sm text-muted-foreground">No account?</Text>
          <Link href="/signup" className="text-sm font-medium text-primary">
            Sign up
          </Link>
        </View>
      </Card>
    </View>
  );
}
