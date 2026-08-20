import type { Meta, StoryObj } from "@storybook/react-vite";

import { Input } from "./input";

const meta = {
  title: "Native/Input",
  component: Input,
  tags: ["autodocs"],
  args: { placeholder: "Type here…", className: "w-72" },
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
export const WithValue: Story = { args: { defaultValue: "hello" } };
