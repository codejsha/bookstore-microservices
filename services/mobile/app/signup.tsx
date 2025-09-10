import { Button } from "@bookstore/design/native/button";
import { Card } from "@bookstore/design/native/card";
import { Link, useRouter } from "expo-router";
import { useEffect } from "react";
import { Text, View } from "react-native";
import { isSessionValid, useAuthStore, useKeycloakAuth } from "@/domains/auth";

export default function SignupScreen() {
  const router = useRouter();
  const authed = useAuthStore(isSessionValid);
  const { request, promptAsync } = useKeycloakAuth({ register: true });

  useEffect(() => {
    if (authed) router.replace("/(tabs)");
  }, [authed, router]);

  return (
    <View className="flex-1 justify-center bg-background px-5">
      <Card className="gap-4">
        <View className="items-center">
          <Text className="text-2xl font-bold text-foreground">
            Create account
          </Text>
          <Text className="mt-1 text-sm text-muted-foreground">
            Sign up with Keycloak.
          </Text>
        </View>

        <Button disabled={!request} onPress={() => promptAsync()}>
          Continue with Keycloak
        </Button>

        <View className="flex-row justify-center gap-2">
          <Text className="text-sm text-muted-foreground">
            Already have an account?
          </Text>
          <Link href="/login" className="text-sm font-medium text-primary">
            Login
          </Link>
        </View>
      </Card>
    </View>
  );
}
