import { Redirect, Tabs } from "expo-router";
import {
  BookOpen,
  ClipboardList,
  Heart,
  Home,
  ShoppingCart,
  User,
} from "lucide-react-native";
import { isSessionValid, useAuthStore } from "@/domains/auth";
import { useThemeColor } from "@/shared/theme/useThemeColor";

export default function TabLayout() {
  const { color } = useThemeColor();
  const authed = useAuthStore(isSessionValid);

  if (!authed) {
    return <Redirect href="/login" />;
  }

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: color("primary"),
        tabBarInactiveTintColor: color("muted-foreground"),
        tabBarStyle: {
          backgroundColor: color("card"),
          borderTopColor: color("border"),
        },
        headerStyle: { backgroundColor: color("card") },
        headerTintColor: color("foreground"),
        sceneStyle: { backgroundColor: color("background") },
        headerShown: true,
      }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: "Home",
          tabBarIcon: ({ color, size }) => <Home color={color} size={size} />,
        }}
      />
      <Tabs.Screen
        name="books"
        options={{
          title: "Books",
          headerShown: false,
          tabBarIcon: ({ color, size }) => (
            <BookOpen color={color} size={size} />
          ),
        }}
      />
      <Tabs.Screen
        name="wishlist"
        options={{
          title: "Wishlist",
          tabBarIcon: ({ color, size }) => <Heart color={color} size={size} />,
        }}
      />
      <Tabs.Screen
        name="cart"
        options={{
          title: "Cart",
          tabBarIcon: ({ color, size }) => (
            <ShoppingCart color={color} size={size} />
          ),
        }}
      />
      <Tabs.Screen
        name="orders"
        options={{
          title: "Orders",
          headerShown: false,
          tabBarIcon: ({ color, size }) => (
            <ClipboardList color={color} size={size} />
          ),
        }}
      />
      <Tabs.Screen
        name="account"
        options={{
          title: "Account",
          tabBarIcon: ({ color, size }) => <User color={color} size={size} />,
        }}
      />
    </Tabs>
  );
}
