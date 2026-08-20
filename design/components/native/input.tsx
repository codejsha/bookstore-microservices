/// <reference path="./nativewind-env.d.ts" />
import { useColorScheme } from "nativewind";
import { TextInput, type TextInputProps } from "react-native";
import { tokens as darkTokens } from "../../build/ts/dark";
import { tokens as lightTokens } from "../../build/ts/light";
import { cn } from "./lib/utils";

interface InputProps extends TextInputProps {
  className?: string;
}

export function Input({
  className,
  placeholderTextColor,
  ...props
}: InputProps) {
  const { colorScheme } = useColorScheme();
  const tokens = colorScheme === "dark" ? darkTokens : lightTokens;

  return (
    <TextInput
      className={cn(
        "h-12 rounded-md border border-input bg-background px-3 text-base text-foreground",
        className,
      )}
      placeholderTextColor={placeholderTextColor ?? tokens["muted-foreground"]}
      {...props}
    />
  );
}
