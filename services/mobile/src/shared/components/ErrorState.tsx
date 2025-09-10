import { Button } from "@bookstore/design/native/button";
import { AlertTriangle } from "lucide-react-native";
import { Text, View } from "react-native";
import { useThemeColor } from "@/shared/theme/useThemeColor";

interface ErrorStateProps {
  message?: string;
  onRetry?: () => void;
  className?: string;
}

export function ErrorState({
  message = "Something went wrong.",
  onRetry,
  className,
}: ErrorStateProps) {
  const { color } = useThemeColor();
  return (
    <View
      className={
        className ?? "flex-1 items-center justify-center bg-background px-8"
      }
    >
      <AlertTriangle color={color("destructive")} size={40} />
      <Text className="mt-3 text-base font-medium text-foreground">
        {message}
      </Text>
      <Text className="mt-1 text-center text-sm text-muted-foreground">
        Please check your connection and try again.
      </Text>
      {onRetry ? (
        <View className="mt-4">
          <Button variant="outline" onPress={onRetry}>
            Retry
          </Button>
        </View>
      ) : null}
    </View>
  );
}
