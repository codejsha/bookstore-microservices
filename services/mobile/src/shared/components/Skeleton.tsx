import { Card } from "@bookstore/design/native/card";
import { useEffect, useRef } from "react";
import { Animated, View } from "react-native";
import { cn } from "@/shared/lib/utils";

export function Skeleton({ className }: { className?: string }) {
  const opacity = useRef(new Animated.Value(0.4)).current;

  useEffect(() => {
    const loop = Animated.loop(
      Animated.sequence([
        Animated.timing(opacity, {
          toValue: 1,
          duration: 700,
          useNativeDriver: true,
        }),
        Animated.timing(opacity, {
          toValue: 0.4,
          duration: 700,
          useNativeDriver: true,
        }),
      ]),
    );
    loop.start();
    return () => loop.stop();
  }, [opacity]);

  return (
    <Animated.View style={{ opacity }}>
      <View className={cn("rounded-md bg-muted", className)} />
    </Animated.View>
  );
}

export function SkeletonCard({ className }: { className?: string }) {
  return (
    <Card className={cn("flex-row items-center gap-3", className)}>
      <Skeleton className="h-20 w-14" />
      <View className="flex-1 gap-2">
        <Skeleton className="h-4 w-3/4" />
        <Skeleton className="h-3 w-1/2" />
        <Skeleton className="h-4 w-16" />
      </View>
    </Card>
  );
}

export function SkeletonList({
  count = 6,
  className,
}: {
  count?: number;
  className?: string;
}) {
  return (
    <View className={cn("gap-3 p-5", className)}>
      {Array.from({ length: count }).map((_, i) => (
        // biome-ignore lint/suspicious/noArrayIndexKey: static placeholder list
        <SkeletonCard key={i} />
      ))}
    </View>
  );
}
