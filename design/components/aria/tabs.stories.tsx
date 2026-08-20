import type { Meta, StoryObj } from "@storybook/react-vite";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "./tabs";

const meta = {
  title: "Aria/Tabs",
  component: Tabs,
  tags: ["autodocs"],
} satisfies Meta<typeof Tabs>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Tabs defaultSelectedKey="account" className="w-80">
      <TabsList aria-label="Settings">
        <TabsTrigger id="account">Account</TabsTrigger>
        <TabsTrigger id="password">Password</TabsTrigger>
        <TabsTrigger id="billing">Billing</TabsTrigger>
      </TabsList>
      <TabsContent id="account">Account settings panel.</TabsContent>
      <TabsContent id="password">Password panel.</TabsContent>
      <TabsContent id="billing">Billing panel.</TabsContent>
    </Tabs>
  ),
};

export const LineVariant: Story = {
  render: () => (
    <Tabs defaultSelectedKey="one" className="w-80">
      <TabsList variant="line" aria-label="Product">
        <TabsTrigger id="one">Overview</TabsTrigger>
        <TabsTrigger id="two">Reviews</TabsTrigger>
        <TabsTrigger id="three">Details</TabsTrigger>
      </TabsList>
      <TabsContent id="one">Overview</TabsContent>
      <TabsContent id="two">Reviews</TabsContent>
      <TabsContent id="three">Details</TabsContent>
    </Tabs>
  ),
};

export const Vertical: Story = {
  render: () => (
    <Tabs defaultSelectedKey="a" orientation="vertical" className="h-48">
      <TabsList aria-label="Nav">
        <TabsTrigger id="a">One</TabsTrigger>
        <TabsTrigger id="b">Two</TabsTrigger>
        <TabsTrigger id="c">Three</TabsTrigger>
      </TabsList>
      <TabsContent id="a">First</TabsContent>
      <TabsContent id="b">Second</TabsContent>
      <TabsContent id="c">Third</TabsContent>
    </Tabs>
  ),
};
