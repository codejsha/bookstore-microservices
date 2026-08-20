import {
  TextArea as AriaTextArea,
  type TextAreaProps as AriaTextAreaProps,
  composeRenderProps,
} from "react-aria-components";

import { cn } from "../lib/utils";

export {
  FieldDescription,
  FieldError,
  Label,
  TextField,
} from "./input";

function Textarea({ className, ...props }: AriaTextAreaProps) {
  return (
    <AriaTextArea
      data-slot="textarea"
      className={composeRenderProps(className, (cls) =>
        cn(
          "flex field-sizing-content min-h-16 w-full rounded-lg border border-input bg-transparent px-2.5 py-2 text-base transition-colors outline-none placeholder:text-muted-foreground data-[focused]:border-ring data-[focused]:ring-3 data-[focused]:ring-ring/50 data-[disabled]:cursor-not-allowed data-[disabled]:bg-input/50 data-[disabled]:opacity-50 data-[invalid]:border-destructive data-[invalid]:ring-3 data-[invalid]:ring-destructive/20 md:text-sm dark:bg-input/30 dark:data-[disabled]:bg-input/80 dark:data-[invalid]:border-destructive/50 dark:data-[invalid]:ring-destructive/40",
          cls,
        ),
      )}
      {...props}
    />
  );
}

export { Textarea };
