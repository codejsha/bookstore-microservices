import type { Meta, StoryObj } from "@storybook/react-vite";

import { Button } from "./button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./card";

const meta = {
  title: "UI/Card",
  component: Card,
  tags: ["autodocs"],
} satisfies Meta<typeof Card>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Basic: Story = {
  render: () => (
    <Card className="w-80">
      <CardHeader>
        <CardTitle>The Pragmatic Programmer</CardTitle>
        <CardDescription>Andrew Hunt, David Thomas</CardDescription>
      </CardHeader>
      <CardContent>
        Your journey to mastery, 20th Anniversary Edition.
      </CardContent>
      <CardFooter>
        <Button size="sm">Add to cart</Button>
      </CardFooter>
    </Card>
  ),
};

export const WithAction: Story = {
  render: () => (
    <Card className="w-80">
      <CardHeader>
        <CardTitle>Wishlist</CardTitle>
        <CardDescription>3 items</CardDescription>
        <CardAction>
          <Button size="icon-sm" variant="ghost">
            ···
          </Button>
        </CardAction>
      </CardHeader>
      <CardContent>You have items saved for later.</CardContent>
    </Card>
  ),
};

export const Small: Story = {
  render: () => (
    <Card size="sm" className="w-72">
      <CardHeader>
        <CardTitle>Compact</CardTitle>
        <CardDescription>Tighter spacing variant.</CardDescription>
      </CardHeader>
      <CardContent>Useful inside dense list layouts.</CardContent>
    </Card>
  ),
};
