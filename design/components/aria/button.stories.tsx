import type { Meta, StoryObj } from "@storybook/react-vite";
import { ArrowRightIcon, PlusIcon } from "lucide-react";

import { Button } from "./button";

const meta = {
  title: "Aria/Button",
  component: Button,
  tags: ["autodocs"],
  argTypes: {
    variant: {
      control: "select",
      options: [
        "default",
        "outline",
        "secondary",
        "ghost",
        "destructive",
        "link",
      ],
    },
    size: {
      control: "select",
      options: [
        "default",
        "xs",
        "sm",
        "lg",
        "icon",
        "icon-xs",
        "icon-sm",
        "icon-lg",
      ],
    },
    isDisabled: { control: "boolean" },
  },
  args: {
    children: "Button",
  },
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Outline: Story = { args: { variant: "outline" } };
export const Secondary: Story = { args: { variant: "secondary" } };
export const Ghost: Story = { args: { variant: "ghost" } };
export const Destructive: Story = { args: { variant: "destructive" } };
export const Link: Story = { args: { variant: "link" } };

export const Disabled: Story = { args: { isDisabled: true } };

export const WithLeadingIcon: Story = {
  args: {
    children: (
      <>
        <PlusIcon />
        Create
      </>
    ),
  },
};

export const WithTrailingIcon: Story = {
  args: {
    variant: "outline",
    children: (
      <>
        Next
        <ArrowRightIcon />
      </>
    ),
  },
};

export const Sizes: Story = {
  render: (args) => (
    <div className="flex items-center gap-3">
      <Button {...args} size="xs">
        XS
      </Button>
      <Button {...args} size="sm">
        SM
      </Button>
      <Button {...args} size="default">
        Default
      </Button>
      <Button {...args} size="lg">
        LG
      </Button>
    </div>
  ),
};
