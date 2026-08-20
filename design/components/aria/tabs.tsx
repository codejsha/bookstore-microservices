import { cva, type VariantProps } from "class-variance-authority";
import {
  Tab as AriaTab,
  TabList as AriaTabList,
  type TabListProps as AriaTabListProps,
  TabPanel as AriaTabPanel,
  type TabPanelProps as AriaTabPanelProps,
  type TabProps as AriaTabProps,
  Tabs as AriaTabs,
  type TabsProps as AriaTabsProps,
  composeRenderProps,
} from "react-aria-components";

import { cn } from "../lib/utils";

function Tabs({
  className,
  orientation = "horizontal",
  ...props
}: AriaTabsProps) {
  return (
    <AriaTabs
      data-slot="tabs"
      orientation={orientation}
      className={composeRenderProps(className, (cls) =>
        cn("group/tabs flex gap-2 data-[orientation=horizontal]:flex-col", cls),
      )}
      {...props}
    />
  );
}

const tabsListVariants = cva(
  "group/tabs-list inline-flex w-fit items-center justify-center rounded-lg p-[3px] text-muted-foreground group-data-[orientation=horizontal]/tabs:h-8 group-data-[orientation=vertical]/tabs:h-fit group-data-[orientation=vertical]/tabs:flex-col data-[variant=line]:rounded-none",
  {
    variants: {
      variant: {
        default: "bg-muted",
        line: "gap-1 bg-transparent",
      },
    },
    defaultVariants: { variant: "default" },
  },
);

function TabsList<T extends object>({
  className,
  variant = "default",
  ...props
}: AriaTabListProps<T> & VariantProps<typeof tabsListVariants>) {
  return (
    <AriaTabList
      data-slot="tabs-list"
      data-variant={variant}
      className={composeRenderProps(className, (cls) =>
        cn(tabsListVariants({ variant }), cls),
      )}
      {...props}
    />
  );
}

function TabsTrigger({ className, ...props }: AriaTabProps) {
  return (
    <AriaTab
      data-slot="tabs-trigger"
      className={composeRenderProps(className, (cls) =>
        cn(
          "relative inline-flex h-[calc(100%-1px)] flex-1 cursor-default items-center justify-center gap-1.5 rounded-md border border-transparent px-1.5 py-0.5 text-sm font-medium whitespace-nowrap text-foreground/60 transition-all group-data-[orientation=vertical]/tabs:w-full group-data-[orientation=vertical]/tabs:justify-start data-[hovered]:text-foreground data-[focus-visible]:border-ring data-[focus-visible]:ring-[3px] data-[focus-visible]:ring-ring/50 data-[disabled]:pointer-events-none data-[disabled]:opacity-50 dark:text-muted-foreground dark:data-[hovered]:text-foreground group-data-[variant=default]/tabs-list:data-[selected]:shadow-sm group-data-[variant=line]/tabs-list:data-[selected]:shadow-none [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 data-[selected]:bg-background data-[selected]:text-foreground dark:data-[selected]:border-input dark:data-[selected]:bg-input/30 group-data-[variant=line]/tabs-list:data-[selected]:bg-transparent dark:group-data-[variant=line]/tabs-list:data-[selected]:bg-transparent",
          cls,
        ),
      )}
      {...props}
    />
  );
}

function TabsContent({ className, ...props }: AriaTabPanelProps) {
  return (
    <AriaTabPanel
      data-slot="tabs-content"
      className={composeRenderProps(className, (cls) =>
        cn("flex-1 text-sm outline-none", cls),
      )}
      {...props}
    />
  );
}

export { Tabs, TabsContent, TabsList, TabsTrigger, tabsListVariants };
