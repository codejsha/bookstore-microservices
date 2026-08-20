import "../global.css";

import { QueryClientProvider } from "@tanstack/react-query";
import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";
import { useColorScheme } from "nativewind";
import { useEffect } from "react";
import { SafeAreaProvider } from "react-native-safe-area-context";

import { queryClient } from "@/shared/api/queryClient";
import { useRouterTelemetry } from "@/shared/hooks/useRouterTelemetry";
import { initTelemetry } from "@/shared/lib/telemetry";
import { useThemeColor } from "@/shared/theme/useThemeColor";
import { ToastProvider } from "@/shared/toast";

function RootNavigator() {
  const { color } = useThemeColor();

  useRouterTelemetry();

  return (
    <Stack
      screenOptions={{
        headerStyle: { backgroundColor: color("card") },
        headerTintColor: color("foreground"),
        contentStyle: { backgroundColor: color("background") },
      }}
    >
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
      <Stack.Screen name="login" options={{ title: "Login" }} />
      <Stack.Screen name="signup" options={{ title: "Sign Up" }} />
    </Stack>
  );
}

export default function RootLayout() {
  const { colorScheme } = useColorScheme();

  useEffect(() => {
    try {
      initTelemetry();
    } catch (e) {
      console.warn("[OTel] Failed to initialize telemetry:", e);
    }
  }, []);

  return (
    <SafeAreaProvider>
      <QueryClientProvider client={queryClient}>
        <ToastProvider>
          <StatusBar style={colorScheme === "dark" ? "light" : "dark"} />
          <RootNavigator />
        </ToastProvider>
      </QueryClientProvider>
    </SafeAreaProvider>
  );
}
