export interface CartItem {
  uid: string;
  product_id: number;
  product_name?: string;
  quantity: number;
  currency: string;
  price: number;
  subtotal: number;
}

export interface Cart {
  uid: string;
  user_id: number;
  items: CartItem[];
  total_items: number;
  total_amount: number;
}

export interface CartAddItemReq {
  product_id: number;
  product_name?: string;
  quantity: number;
  currency: string;
  price: number;
}

export interface CartUpdateItemReq {
  quantity: number;
}

export interface CartCheckoutReq {
  currency: string;
  idempotency_key: string;
}
