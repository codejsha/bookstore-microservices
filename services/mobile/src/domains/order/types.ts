export type OrderStatus =
  | "PENDING"
  | "PAID"
  | "SHIPPED"
  | "DELIVERED"
  | "CANCELLED"
  | "REFUNDED";

export interface OrderItem {
  uid: string;
  productId: number;
  productName?: string;
  sku?: string;
  quantity: number;
  currency: string;
  price: number;
  taxRate?: number;
  subtotal: number;
  createdAt?: string;
}

export interface OrderShipping {
  uid: string;
  recipientName: string;
  recipientPhone?: string;
  addressLine1: string;
  addressLine2?: string;
  city: string;
  state?: string;
  postalCode?: string;
  country: string;
  shippingMethod?: string;
  createdAt?: string;
}

export interface Order {
  uid: string;
  userId: number;
  orderNumber: string;
  status: OrderStatus;
  currency: string;
  itemsAmount: number;
  discountAmount: number;
  shippingAmount: number;
  taxAmount: number;
  totalAmount: number;
  items?: OrderItem[];
  shipping?: OrderShipping;
  createdAt: string;
  updatedAt?: string;
}

export interface OrderFindAllResp {
  total: number;
  items: Order[];
}

export interface OrderQueryParams {
  user_id?: number;
  status?: string;
  size?: number;
  page?: number;
  sort?: string;
}
