import type { Meta, StoryObj } from "@storybook/react-vite";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "./tabs";

const meta = {
  title: "UI/Tabs",
  component: Tabs,
  tags: ["autodocs"],
} satisfies Meta<typeof Tabs>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Tabs defaultValue="account" className="w-80">
      <TabsList>
        <TabsTrigger value="account">Account</TabsTrigger>
        <TabsTrigger value="password">Password</TabsTrigger>
        <TabsTrigger value="billing">Billing</TabsTrigger>
      </TabsList>
      <TabsContent value="account">Account settings panel.</TabsContent>
      <TabsContent value="password">Password panel.</TabsContent>
      <TabsContent value="billing">Billing panel.</TabsContent>
    </Tabs>
  ),
};

export const LineVariant: Story = {
  render: () => (
    <Tabs defaultValue="one" className="w-80">
      <TabsList variant="line">
        <TabsTrigger value="one">Overview</TabsTrigger>
        <TabsTrigger value="two">Reviews</TabsTrigger>
        <TabsTrigger value="three">Details</TabsTrigger>
      </TabsList>
      <TabsContent value="one">Overview content</TabsContent>
      <TabsContent value="two">Reviews content</TabsContent>
      <TabsContent value="three">Details content</TabsContent>
    </Tabs>
  ),
};

export const Vertical: Story = {
  render: () => (
    <Tabs defaultValue="a" orientation="vertical" className="h-48">
      <TabsList>
        <TabsTrigger value="a">One</TabsTrigger>
        <TabsTrigger value="b">Two</TabsTrigger>
        <TabsTrigger value="c">Three</TabsTrigger>
      </TabsList>
      <TabsContent value="a">First</TabsContent>
      <TabsContent value="b">Second</TabsContent>
      <TabsContent value="c">Third</TabsContent>
    </Tabs>
  ),
};
