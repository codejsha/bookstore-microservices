import type { Meta, StoryObj } from "@storybook/react-vite";

import { Label } from "./label";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectSeparator,
  SelectTrigger,
  SelectValue,
} from "./select";

const meta = {
  title: "Aria/Select",
  component: Select,
  tags: ["autodocs"],
} satisfies Meta<typeof Select>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Basic: Story = {
  render: () => (
    <Select className="flex flex-col gap-1.5">
      <Label>Fruit</Label>
      <SelectTrigger className="w-48">
        <SelectValue placeholder="Pick a fruit" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem id="apple">Apple</SelectItem>
        <SelectItem id="banana">Banana</SelectItem>
        <SelectItem id="cherry">Cherry</SelectItem>
        <SelectItem id="durian">Durian</SelectItem>
      </SelectContent>
    </Select>
  ),
};

export const WithGroups: Story = {
  render: () => (
    <Select className="flex flex-col gap-1.5">
      <Label>Category</Label>
      <SelectTrigger className="w-48">
        <SelectValue placeholder="Select a category" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Fiction</SelectLabel>
          <SelectItem id="novel">Novel</SelectItem>
          <SelectItem id="poetry">Poetry</SelectItem>
        </SelectGroup>
        <SelectSeparator />
        <SelectGroup>
          <SelectLabel>Non-Fiction</SelectLabel>
          <SelectItem id="biography">Biography</SelectItem>
          <SelectItem id="history">History</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  ),
};

export const Disabled: Story = {
  render: () => (
    <Select isDisabled className="flex flex-col gap-1.5">
      <Label>Locale</Label>
      <SelectTrigger className="w-48">
        <SelectValue placeholder="Disabled" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem id="en">English</SelectItem>
      </SelectContent>
    </Select>
  ),
};

export const Small: Story = {
  render: () => (
    <Select className="flex flex-col gap-1.5">
      <Label>Size</Label>
      <SelectTrigger size="sm" className="w-40">
        <SelectValue placeholder="Small" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem id="a">Option A</SelectItem>
        <SelectItem id="b">Option B</SelectItem>
      </SelectContent>
    </Select>
  ),
};
