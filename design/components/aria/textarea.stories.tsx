import type { Meta, StoryObj } from "@storybook/react-vite";

import { Label, Textarea, TextField } from "./textarea";

const meta = {
  title: "Aria/Textarea",
  component: Textarea,
  tags: ["autodocs"],
  args: { placeholder: "Share your thoughts…" },
} satisfies Meta<typeof Textarea>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithLabel: Story = {
  render: (args) => (
    <TextField className="w-80">
      <Label>Bio</Label>
      <Textarea {...args} />
    </TextField>
  ),
};

export const Disabled: Story = {
  render: () => (
    <TextField isDisabled defaultValue="read-only" className="w-80">
      <Label>Notes</Label>
      <Textarea />
    </TextField>
  ),
};

export const Invalid: Story = {
  render: () => (
    <TextField isInvalid defaultValue="too short" className="w-80">
      <Label>Description</Label>
      <Textarea />
    </TextField>
  ),
};
