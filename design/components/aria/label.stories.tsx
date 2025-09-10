import type { Meta, StoryObj } from "@storybook/react-vite";

import { Input, Label, TextField } from "./input";

const meta = {
  title: "Aria/Label",
  component: Label,
  tags: ["autodocs"],
  args: { children: "Email" },
} satisfies Meta<typeof Label>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithInput: Story = {
  render: () => (
    <TextField className="w-72">
      <Label>Email</Label>
      <Input type="email" placeholder="you@example.com" />
    </TextField>
  ),
};
