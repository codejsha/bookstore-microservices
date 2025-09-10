import { CheckIcon, ChevronDownIcon } from "lucide-react";
import {
  Header as AriaHeader,
  type HeaderProps as AriaHeaderProps,
  ListBox as AriaListBox,
  ListBoxItem as AriaListBoxItem,
  type ListBoxItemProps as AriaListBoxItemProps,
  type ListBoxProps as AriaListBoxProps,
  ListBoxSection as AriaListBoxSection,
  type ListBoxSectionProps as AriaListBoxSectionProps,
  Popover as AriaPopover,
  type PopoverProps as AriaPopoverProps,
  Select as AriaSelect,
  type SelectProps as AriaSelectProps,
  SelectValue as AriaSelectValue,
  type SelectValueProps as AriaSelectValueProps,
  Separator as AriaSeparator,
  composeRenderProps,
} from "react-aria-components";
import { cn } from "../lib/utils";
import { Button } from "./button";

function Select<T extends object>(props: AriaSelectProps<T>) {
  return <AriaSelect data-slot="select" {...props} />;
}

function SelectValue<T extends object>({
  className,
  placeholder,
  ...props
}: AriaSelectValueProps<T> & { placeholder?: string }) {
  return (
    <AriaSelectValue
      data-slot="select-value"
      className={composeRenderProps(className, (cls) =>
        cn(
          "flex flex-1 text-left data-[placeholder]:text-muted-foreground",
          cls,
        ),
      )}
      {...props}
    >
      {composeRenderProps(props.children, (resolved, renderProps) =>
        renderProps.isPlaceholder && placeholder
          ? placeholder
          : (resolved ?? renderProps.defaultChildren),
      )}
    </AriaSelectValue>
  );
}

function SelectTrigger({
  className,
  size = "default",
  children,
  ...props
}: React.ComponentProps<typeof Button> & { size?: "default" | "sm" }) {
  return (
    <Button
      data-slot="select-trigger"
      data-size={size}
      className={composeRenderProps(className, (cls) =>
        cn(
          "flex w-fit items-center justify-between gap-1.5 rounded-lg border border-input bg-transparent py-2 pr-2 pl-2.5 text-sm font-normal whitespace-nowrap data-[focus-visible]:border-ring data-[focus-visible]:ring-3 data-[focus-visible]:ring-ring/50 data-[disabled]:cursor-not-allowed data-[disabled]:opacity-50 data-[size=default]:h-8 data-[size=sm]:h-7 data-[size=sm]:rounded-[min(var(--radius-md),10px)] dark:bg-input/30 dark:data-[hovered]:bg-input/50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
          cls,
        ),
      )}
      {...props}
    >
      {composeRenderProps(children, (resolved) => (
        <>
          {resolved}
          <ChevronDownIcon className="pointer-events-none size-4 text-muted-foreground" />
        </>
      ))}
    </Button>
  );
}

function SelectPopover({ className, ...props }: AriaPopoverProps) {
  return (
    <AriaPopover
      data-slot="select-popover"
      className={composeRenderProps(className, (cls) =>
        cn(
          "z-50 max-h-(--trigger-height) w-(--trigger-width) min-w-36 origin-(--origin) overflow-x-hidden overflow-y-auto rounded-lg bg-popover text-popover-foreground shadow-md ring-1 ring-foreground/10 outline-none data-[entering]:animate-in data-[entering]:fade-in-0 data-[entering]:zoom-in-95 data-[exiting]:animate-out data-[exiting]:fade-out-0 data-[exiting]:zoom-out-95 data-[placement=bottom]:slide-in-from-top-2 data-[placement=top]:slide-in-from-bottom-2",
          cls,
        ),
      )}
      {...props}
    />
  );
}

function SelectContent<T extends object>({
  className,
  ...props
}: AriaListBoxProps<T>) {
  return (
    <SelectPopover>
      <AriaListBox
        data-slot="select-content"
        className={composeRenderProps(className, (cls) =>
          cn("p-1 outline-none", cls),
        )}
        {...props}
      />
    </SelectPopover>
  );
}

function SelectLabel({ className, ...props }: AriaHeaderProps) {
  return (
    <AriaHeader
      data-slot="select-label"
      className={cn("px-1.5 py-1 text-xs text-muted-foreground", className)}
      {...props}
    />
  );
}

function SelectItem({ className, children, ...props }: AriaListBoxItemProps) {
  return (
    <AriaListBoxItem
      data-slot="select-item"
      className={composeRenderProps(className, (cls) =>
        cn(
          "relative flex w-full cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none data-[focused]:bg-accent data-[focused]:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
          cls,
        ),
      )}
      {...props}
    >
      {composeRenderProps(children, (resolved, { isSelected }) => (
        <>
          {resolved}
          {isSelected && (
            <span className="pointer-events-none absolute right-2 flex size-4 items-center justify-center">
              <CheckIcon />
            </span>
          )}
        </>
      ))}
    </AriaListBoxItem>
  );
}

function SelectGroup<T extends object>({
  className,
  ...props
}: AriaListBoxSectionProps<T>) {
  return (
    <AriaListBoxSection
      data-slot="select-group"
      className={cn(className)}
      {...props}
    />
  );
}

function SelectSeparator({
  className,
  ...props
}: React.ComponentProps<typeof AriaSeparator>) {
  return (
    <AriaSeparator
      data-slot="select-separator"
      className={cn("pointer-events-none -mx-1 my-1 h-px bg-border", className)}
      {...props}
    />
  );
}

export {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectPopover,
  SelectSeparator,
  SelectTrigger,
  SelectValue,
};
