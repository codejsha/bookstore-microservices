import type { Meta, StoryObj } from "@storybook/react-vite";

import { Button } from "./button";
import { Card, CardContent, CardHeader, CardTitle } from "./card";

const meta = {
  title: "Native/Card",
  component: Card,
  tags: ["autodocs"],
} satisfies Meta<typeof Card>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Basic: Story = {
  args: {
    className: "w-80",
    children: (
      <>
        <CardHeader>
          <CardTitle>React Native</CardTitle>
        </CardHeader>
        <CardContent>
          <Button>Add to cart</Button>
        </CardContent>
      </>
    ),
  },
};
