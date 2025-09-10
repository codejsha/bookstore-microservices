import {
  Label as AriaLabel,
  type LabelProps as AriaLabelProps,
} from "react-aria-components";

import { cn } from "../lib/utils";

function Label({ className, ...props }: AriaLabelProps) {
  return (
    <AriaLabel
      data-slot="label"
      className={cn(
        "flex items-center gap-2 text-sm leading-none font-medium select-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
        className,
      )}
      {...props}
    />
  );
}

export { Label };
