import {
  FieldError as AriaFieldError,
  type FieldErrorProps as AriaFieldErrorProps,
  Input as AriaInput,
  type InputProps as AriaInputProps,
  TextField as AriaTextField,
  type TextFieldProps as AriaTextFieldProps,
  type TextProps as AriaTextProps,
  composeRenderProps,
  Text,
} from "react-aria-components";
import { cn } from "../lib/utils";
import { Label } from "./label";

function TextField({ className, ...props }: AriaTextFieldProps) {
  return (
    <AriaTextField
      data-slot="text-field"
      className={composeRenderProps(className, (cls) =>
        cn("flex flex-col gap-1.5", cls),
      )}
      {...props}
    />
  );
}

function Input({ className, ...props }: AriaInputProps) {
  return (
    <AriaInput
      data-slot="input"
      className={composeRenderProps(className, (cls) =>
        cn(
          "h-8 w-full min-w-0 rounded-lg border border-input bg-transparent px-2.5 py-1 text-base transition-colors outline-none placeholder:text-muted-foreground data-[focused]:border-ring data-[focused]:ring-3 data-[focused]:ring-ring/50 data-[disabled]:pointer-events-none data-[disabled]:cursor-not-allowed data-[disabled]:bg-input/50 data-[disabled]:opacity-50 data-[invalid]:border-destructive data-[invalid]:ring-3 data-[invalid]:ring-destructive/20 md:text-sm dark:bg-input/30 dark:data-[disabled]:bg-input/80 dark:data-[invalid]:border-destructive/50 dark:data-[invalid]:ring-destructive/40",
          cls,
        ),
      )}
      {...props}
    />
  );
}

function FieldDescription({ className, ...props }: AriaTextProps) {
  return (
    <Text
      slot="description"
      data-slot="field-description"
      className={cn("text-xs text-muted-foreground", className)}
      {...props}
    />
  );
}

function FieldError({ className, ...props }: AriaFieldErrorProps) {
  return (
    <AriaFieldError
      data-slot="field-error"
      className={composeRenderProps(className, (cls) =>
        cn("text-xs text-destructive", cls),
      )}
      {...props}
    />
  );
}

export { FieldDescription, FieldError, Input, Label, TextField };
