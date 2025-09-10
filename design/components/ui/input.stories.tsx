import type { Meta, StoryObj } from "@storybook/react-vite";

import { Input } from "./input";
import { Label } from "./label";

const meta = {
  title: "UI/Input",
  component: Input,
  tags: ["autodocs"],
  args: { placeholder: "Type here…" },
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithLabel: Story = {
  render: (args) => (
    <div className="flex w-72 flex-col gap-1.5">
      <Label htmlFor="email">Email</Label>
      <Input id="email" type="email" {...args} />
    </div>
  ),
};

export const Disabled: Story = {
  args: { disabled: true, value: "Read-only" },
};

export const Invalid: Story = {
  args: { "aria-invalid": true, defaultValue: "not-an-email" },
};

export const Required: Story = {
  render: (args) => (
    <div className="flex w-72 flex-col gap-1.5">
      <Label htmlFor="full-name">Full name</Label>
      <Input id="full-name" required placeholder="Jane Doe" {...args} />
    </div>
  ),
};

export const Types: Story = {
  render: () => (
    <div className="flex w-72 flex-col gap-3">
      <Input placeholder="text" />
      <Input type="email" placeholder="email" />
      <Input type="password" placeholder="password" />
      <Input type="number" placeholder="number" />
      <Input type="file" />
    </div>
  ),
};
