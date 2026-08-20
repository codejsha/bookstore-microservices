/// <reference path="./nativewind-env.d.ts" />
import { cva, type VariantProps } from "class-variance-authority";
import { Pressable, Text } from "react-native";
import { cn } from "./lib/utils";

const buttonVariants = cva("flex-row items-center justify-center rounded-md", {
  variants: {
    variant: {
      default: "bg-primary",
      destructive: "bg-destructive",
      outline: "border border-border bg-background",
      secondary: "bg-secondary",
      ghost: "",
      link: "",
    },
    size: {
      default: "h-12 px-4 py-2",
      sm: "h-9 px-3",
      lg: "h-14 px-8",
      icon: "h-12 w-12",
    },
  },
  defaultVariants: {
    variant: "default",
    size: "default",
  },
});

const buttonTextVariants = cva("text-base font-medium", {
  variants: {
    variant: {
      default: "text-primary-foreground",
      destructive: "text-destructive-foreground",
      outline: "text-foreground",
      secondary: "text-secondary-foreground",
      ghost: "text-foreground",
      link: "text-primary underline",
    },
  },
  defaultVariants: {
    variant: "default",
  },
});

interface ButtonProps extends VariantProps<typeof buttonVariants> {
  className?: string;
  textClassName?: string;
  children: string;
  onPress?: () => void;
  disabled?: boolean;
}

export function Button({
  className,
  textClassName,
  variant,
  size,
  children,
  onPress,
  disabled,
}: ButtonProps) {
  return (
    <Pressable
      className={cn(
        buttonVariants({ variant, size }),
        disabled && "opacity-50",
        className,
      )}
      onPress={onPress}
      disabled={disabled}
    >
      <Text className={cn(buttonTextVariants({ variant }), textClassName)}>
        {children}
      </Text>
    </Pressable>
  );
}
