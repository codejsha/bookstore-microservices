import type { Meta, StoryObj } from "@storybook/react-vite";
import { LogOutIcon, SettingsIcon, UserIcon } from "lucide-react";

import { Button } from "./button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
} from "./dropdown-menu";

const meta = {
  title: "Aria/DropdownMenu",
  component: DropdownMenu,
  tags: ["autodocs"],
  args: { children: null },
} satisfies Meta<typeof DropdownMenu>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Basic: Story = {
  render: () => (
    <DropdownMenu>
      <Button variant="outline">Open menu</Button>
      <DropdownMenuContent>
        <DropdownMenuLabel>My Account</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem>
          <UserIcon />
          Profile
          <DropdownMenuShortcut>⇧⌘P</DropdownMenuShortcut>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <SettingsIcon />
          Settings
          <DropdownMenuShortcut>⌘,</DropdownMenuShortcut>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem variant="destructive">
          <LogOutIcon />
          Sign out
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  ),
};

export const WithCheckbox: Story = {
  render: () => (
    <DropdownMenu>
      <Button variant="outline">View options</Button>
      <DropdownMenuContent
        selectionMode="multiple"
        defaultSelectedKeys={["name", "status"]}
      >
        <DropdownMenuLabel>Columns</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuCheckboxItem id="name">Name</DropdownMenuCheckboxItem>
        <DropdownMenuCheckboxItem id="status">Status</DropdownMenuCheckboxItem>
        <DropdownMenuCheckboxItem id="created">
          Created
        </DropdownMenuCheckboxItem>
      </DropdownMenuContent>
    </DropdownMenu>
  ),
};

export const WithRadio: Story = {
  render: () => (
    <DropdownMenu>
      <Button variant="outline">Sort</Button>
      <DropdownMenuContent
        selectionMode="single"
        defaultSelectedKeys={["newest"]}
      >
        <DropdownMenuLabel>Sort by</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuRadioItem id="newest">Newest</DropdownMenuRadioItem>
        <DropdownMenuRadioItem id="oldest">Oldest</DropdownMenuRadioItem>
        <DropdownMenuRadioItem id="popular">Popular</DropdownMenuRadioItem>
      </DropdownMenuContent>
    </DropdownMenu>
  ),
};

export const Disabled: Story = {
  render: () => (
    <DropdownMenu>
      <Button variant="outline">Actions</Button>
      <DropdownMenuContent>
        <DropdownMenuItem>Run</DropdownMenuItem>
        <DropdownMenuItem isDisabled>Pause (unavailable)</DropdownMenuItem>
        <DropdownMenuItem>Stop</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  ),
};
