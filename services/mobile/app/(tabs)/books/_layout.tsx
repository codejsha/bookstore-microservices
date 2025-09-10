import { Stack } from "expo-router";

export default function BooksLayout() {
  return (
    <Stack>
      <Stack.Screen name="index" options={{ title: "Books" }} />
      <Stack.Screen name="[uid]" options={{ title: "Book Detail" }} />
    </Stack>
  );
}
