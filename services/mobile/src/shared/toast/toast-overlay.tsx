import { CheckCircle2, Info, XCircle } from "lucide-react-native";
import { useEffect, useRef } from "react";
import { Animated, Pressable, Text, View } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useThemeColor } from "@/shared/theme/useThemeColor";
import type { ToastItem, ToastVariant } from "./toast-context";

const ICONS = {
  success: CheckCircle2,
  error: XCircle,
  info: Info,
} as const;

function iconColorName(variant: ToastVariant) {
  if (variant === "error") return "destructive" as const;
  if (variant === "success") return "chart-1" as const;
  return "primary" as const;
}

function AnimatedToast({
  toast,
  onDismiss,
}: {
  toast: ToastItem;
  onDismiss: (id: number) => void;
}) {
  const { color } = useThemeColor();
  const anim = useRef(new Animated.Value(0)).current;
  const Icon = ICONS[toast.variant];

  useEffect(() => {
    Animated.spring(anim, {
      toValue: 1,
      useNativeDriver: true,
      friction: 8,
      tension: 80,
    }).start();
  }, [anim]);

  return (
    <Animated.View
      style={{
        opacity: anim,
        transform: [
          {
            translateY: anim.interpolate({
              inputRange: [0, 1],
              outputRange: [-16, 0],
            }),
          },
        ],
      }}
    >
      <Pressable onPress={() => onDismiss(toast.id)}>
        <View className="flex-row items-center gap-3 rounded-md border border-border bg-card px-4 py-3 shadow-lg">
          <Icon color={color(iconColorName(toast.variant))} size={20} />
          <Text
            className="flex-1 text-sm font-medium text-card-foreground"
            numberOfLines={3}
          >
            {toast.message}
          </Text>
        </View>
      </Pressable>
    </Animated.View>
  );
}

export function ToastOverlay({
  toasts,
  onDismiss,
}: {
  toasts: ToastItem[];
  onDismiss: (id: number) => void;
}) {
  const insets = useSafeAreaInsets();
  if (toasts.length === 0) return null;

  return (
    <View
      pointerEvents="box-none"
      style={{ top: insets.top + 8 }}
      className="absolute left-0 right-0 z-50 gap-2 px-4"
    >
      {toasts.map((toast) => (
        <AnimatedToast key={toast.id} toast={toast} onDismiss={onDismiss} />
      ))}
    </View>
  );
}
