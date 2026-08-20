/// <reference path="./nativewind-env.d.ts" />
import {Text, View} from "react-native";
import {cn} from "./lib/utils";

interface CardProps {
    className?: string;
    children: React.ReactNode;
}

export function Card({className, children}: CardProps) {
    return (
            <View
                    className={cn("rounded-lg border border-border bg-card p-4", className)}
            >
                {children}
            </View>
    );
}

interface CardHeaderProps {
    className?: string;
    children: React.ReactNode;
}

export function CardHeader({className, children}: CardHeaderProps) {
    return <View className={cn("pb-2", className)}>{children}</View>;
}

interface CardTitleProps {
    className?: string;
    children: string;
}

export function CardTitle({className, children}: CardTitleProps) {
    return (
            <Text
                    className={cn("text-lg font-semibold text-card-foreground", className)}
            >
                {children}
            </Text>
    );
}

interface CardContentProps {
    className?: string;
    children: React.ReactNode;
}

export function CardContent({className, children}: CardContentProps) {
    return <View className={cn("pt-0", className)}>{children}</View>;
}
