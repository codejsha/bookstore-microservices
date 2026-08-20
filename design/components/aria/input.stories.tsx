import type { Meta, StoryObj } from "@storybook/react-vite";

import { FieldDescription, FieldError, Input, Label, TextField } from "./input";

const meta = {
  title: "Aria/Input",
  component: Input,
  tags: ["autodocs"],
  args: { placeholder: "Type here…" },
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithLabel: Story = {
  render: (args) => (
    <TextField className="w-72">
      <Label>Email</Label>
      <Input type="email" {...args} />
      <FieldDescription>We never share your email.</FieldDescription>
    </TextField>
  ),
};

export const Disabled: Story = {
  render: () => (
    <TextField isDisabled defaultValue="read-only" className="w-72">
      <Label>Username</Label>
      <Input />
    </TextField>
  ),
};

export const Invalid: Story = {
  render: () => (
    <TextField isInvalid defaultValue="oops" className="w-72">
      <Label>Token</Label>
      <Input />
      <FieldError>Token format is invalid.</FieldError>
    </TextField>
  ),
};

export const Required: Story = {
  render: () => (
    <TextField isRequired className="w-72" defaultValue="">
      <Label>Full name</Label>
      <Input placeholder="Jane Doe" />
      <FieldError />
    </TextField>
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
