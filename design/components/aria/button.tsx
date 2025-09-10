import { cva, type VariantProps } from "class-variance-authority";
import {
  Button as AriaButton,
  type ButtonProps as AriaButtonProps,
  composeRenderProps,
} from "react-aria-components";

import { cn } from "../lib/utils";

const buttonVariants = cva(
  "inline-flex shrink-0 cursor-default items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none data-[focus-visible]:border-ring data-[focus-visible]:ring-3 data-[focus-visible]:ring-ring/50 data-[pressed]:translate-y-px data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
  {
    variants: {
      variant: {
        default:
          "bg-primary text-primary-foreground data-[hovered]:bg-primary/80",
        outline:
          "border-border bg-background data-[hovered]:bg-muted data-[hovered]:text-foreground dark:border-input dark:bg-input/30 dark:data-[hovered]:bg-input/50",
        secondary:
          "bg-secondary text-secondary-foreground data-[hovered]:bg-secondary/80",
        ghost:
          "data-[hovered]:bg-muted data-[hovered]:text-foreground dark:data-[hovered]:bg-muted/50",
        destructive:
          "bg-destructive/10 text-destructive data-[hovered]:bg-destructive/20 data-[focus-visible]:border-destructive/40 data-[focus-visible]:ring-destructive/20 dark:bg-destructive/20 dark:data-[hovered]:bg-destructive/30 dark:data-[focus-visible]:ring-destructive/40",
        link: "text-primary underline-offset-4 data-[hovered]:underline",
      },
      size: {
        default: "h-8 gap-1.5 px-2.5",
        xs: "h-6 gap-1 rounded-[min(var(--radius-md),10px)] px-2 text-xs [&_svg:not([class*='size-'])]:size-3",
        sm: "h-7 gap-1 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] [&_svg:not([class*='size-'])]:size-3.5",
        lg: "h-9 gap-1.5 px-2.5",
        icon: "size-8",
        "icon-xs":
          "size-6 rounded-[min(var(--radius-md),10px)] [&_svg:not([class*='size-'])]:size-3",
        "icon-sm": "size-7 rounded-[min(var(--radius-md),12px)]",
        "icon-lg": "size-9",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

interface ButtonProps
  extends AriaButtonProps,
    VariantProps<typeof buttonVariants> {}

function Button({ className, variant, size, ...props }: ButtonProps) {
  return (
    <AriaButton
      data-slot="button"
      className={composeRenderProps(className, (cls) =>
        cn(buttonVariants({ variant, size }), cls),
      )}
      {...props}
    />
  );
}

export type { ButtonProps };
export { Button, buttonVariants };
