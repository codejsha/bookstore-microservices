import type { AdminIdentity, Dashboard } from "@/domains/admin";
import type { Author, Subject, Work } from "@/domains/catalog";
import type { Stock, Warehouse } from "@/domains/inventory";
import {
  ORDER_STATUSES,
  type Order,
  type OrderLine,
  type OrderStatus,
} from "@/domains/order";
import {
  PAYMENT_STATUSES,
  type Payment,
  REFUND_STATUSES,
  type Refund,
} from "@/domains/payment";
import type { Settlement, SettlementDetailLine } from "@/domains/settlement";
import type { User, UserStatus } from "@/domains/user";

export const identity: AdminIdentity = {
  uid: "11111111-1111-1111-1111-111111111111",
  email: "admin@example.com",
  name: "root",
  roles: ["MANAGER", "STAFF", "USER"],
};

export const dashboard: Dashboard = {
  works: { count: "1284", available: true },
  orders: { count: "0", available: false },
  warehouses: { count: "3", available: true },
};

export const authors: Author[] = [
  { uid: "33333333-3333-3333-3333-333333333331", name: "Frank Herbert" },
  { uid: "33333333-3333-3333-3333-333333333332", name: "Ursula K. Le Guin" },
  { uid: "33333333-3333-3333-3333-333333333333", name: "Octavia E. Butler" },
];

export const subjects: Subject[] = [
  { uid: "44444444-4444-4444-4444-444444444441", name: "Science Fiction" },
  { uid: "44444444-4444-4444-4444-444444444442", name: "Fantasy" },
];

export const works: Work[] = Array.from({ length: 43 }, (_, i) => ({
  uid: `22222222-2222-2222-2222-${String(i).padStart(12, "0")}`,
  title: `Work ${i + 1}`,
  description: i % 3 === 0 ? `A description for work ${i + 1}.` : undefined,
  first_publish_date: String(1960 + (i % 40)),
  ol_key: i % 2 === 0 ? `OL${i}W` : undefined,
  authors: [authors[i % authors.length]],
  subjects: [subjects[i % subjects.length]],
  created_at: "2026-01-01T00:00:00Z",
  updated_at: undefined,
}));

export const riskEntries = [
  {
    user_uid: "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    level: "restrict",
    reason: "4xx burst on catalog endpoints",
    flagged_by: "11111111-1111-1111-1111-111111111111",
    flagged_at: "2026-08-19T02:10:00Z",
    expires_at: "2026-08-21T02:10:00Z",
  },
  {
    user_uid: "550e8400-e29b-41d4-a716-446655440000",
    level: "block",
    reason: "credential stuffing from shared IP",
    flagged_by: "11111111-1111-1111-1111-111111111111",
    flagged_at: "2026-08-19T09:30:00Z",
    expires_at: "2026-08-26T09:30:00Z",
  },
];

const pad = (n: number) => String(n).padStart(12, "0");
const at = (day: number, hour = 9) =>
  new Date(Date.UTC(2026, 8, day, hour, (day * 7) % 60)).toISOString();

const FIRST_NAMES = ["Ada", "Grace", "Alan", "Barbara", "Edsger", "Margaret"];
const LAST_NAMES = ["Lovelace", "Hopper", "Turing", "Liskov", "Dijkstra"];

const userStatus = (i: number): UserStatus =>
  i % 11 === 7 ? "SUSPENDED" : i % 13 === 12 ? "DEACTIVATED" : "ACTIVE";

export const users: User[] = Array.from({ length: 32 }, (_, i) => {
  const first = FIRST_NAMES[i % FIRST_NAMES.length];
  const last = LAST_NAMES[i % LAST_NAMES.length];
  return {
    uid: `55555555-5555-4555-8555-${pad(i)}`,
    email: `${first.toLowerCase()}.${last.toLowerCase()}${i}@example.com`,
    first_name: first,
    last_name: last,
    phone: i % 3 === 0 ? undefined : `+1-555-01${String(i).padStart(2, "0")}`,
    status: userStatus(i),
    roles: i % 8 === 0 ? ["USER", "STAFF"] : ["USER"],
    last_login_at: i % 4 === 3 ? undefined : at(1 + (i % 28), 14),
    created_at: at(1 + (i % 20)),
    updated_at: i % 5 === 0 ? at(25) : undefined,
  };
});

const PRODUCTS = [
  { name: "Dune", price: { USD: 18.99, KRW: 21000 } },
  { name: "The Left Hand of Darkness", price: { USD: 16.5, KRW: 18500 } },
  { name: "Kindred", price: { USD: 14.25, KRW: 16000 } },
  { name: "The Dispossessed", price: { USD: 17, KRW: 19000 } },
];

const round2 = (n: number) => Math.round(n * 100) / 100;

export const orders: Order[] = Array.from({ length: 45 }, (_, i) => {
  const currency = i % 4 === 0 ? "KRW" : "USD";
  const user = users[i % 12];
  const lineCount = 1 + (i % 3);
  const lines: OrderLine[] = Array.from({ length: lineCount }, (_, j) => {
    const product = PRODUCTS[(i + j) % PRODUCTS.length];
    const quantity = 1 + ((i + j) % 2);
    const price = product.price[currency];
    return {
      uid: `66666666-6666-4666-8666-${pad(i * 10 + j)}`,
      product_id: 1000 + ((i + j) % PRODUCTS.length),
      product_name: product.name,
      sku: `SKU-${1000 + ((i + j) % PRODUCTS.length)}`,
      quantity,
      price,
      subtotal: round2(price * quantity),
      tax_rate: 0.1,
      currency,
      options: j === 1 ? "hardcover" : undefined,
    };
  });
  const itemsAmount = round2(lines.reduce((sum, l) => sum + l.subtotal, 0));
  const discount = i % 5 === 0 ? round2(itemsAmount * 0.1) : 0;
  const shippingAmount = currency === "KRW" ? 3000 : 4.99;
  const tax = round2((itemsAmount - discount) * 0.1);
  const status: OrderStatus = ORDER_STATUSES[i % ORDER_STATUSES.length];
  return {
    uid: `77777777-7777-4777-8777-${pad(i)}`,
    order_number: `ord_${String(9_000_000 + i * 137).padStart(24, "0")}`,
    user_uid: user.uid,
    status,
    currency,
    items_amount: itemsAmount,
    discount_amount: discount,
    shipping_amount: shippingAmount,
    tax_amount: tax,
    total_amount: round2(itemsAmount - discount + shippingAmount + tax),
    items: lines,
    shipping: {
      recipient_name: `${user.first_name} ${user.last_name}`,
      recipient_phone: user.phone ?? "+1-555-0100",
      address_line1: `${100 + i} Market Street`,
      address_line2: i % 2 === 0 ? `Apt ${i + 1}` : undefined,
      city: currency === "KRW" ? "Seoul" : "San Francisco",
      state: currency === "KRW" ? "Seoul" : "CA",
      postal_code: currency === "KRW" ? "04524" : "94103",
      country: currency === "KRW" ? "KR" : "US",
      shipping_method: i % 3 === 0 ? "express" : "standard",
    },
    created_at: at(1 + (i % 28), 10),
    updated_at: status === "PENDING" ? undefined : at(1 + (i % 28), 16),
  };
});

const CONNECTORS = ["stripe", "adyen"];
const PAYMENT_METHODS = ["card", "wallet", "bank_transfer", "pay_later"];

export const payments: Payment[] = Array.from({ length: 30 }, (_, i) => {
  const currency = i % 4 === 0 ? "KRW" : "USD";
  const status = PAYMENT_STATUSES[i % PAYMENT_STATUSES.length];
  const amount = currency === "KRW" ? 15_000 + i * 1_000 : 1_999 + i * 250;
  const captured = status === "succeeded" || status === "partially_captured";
  return {
    uid: `88888888-8888-4888-8888-${pad(i)}`,
    payment_id: `pay_${String(i).padStart(20, "0")}`,
    customer_id:
      i % 6 === 5 ? undefined : `cus_${String(i % 8).padStart(4, "0")}`,
    status,
    amount,
    amount_captured: captured
      ? status === "partially_captured"
        ? Math.floor(amount / 2)
        : amount
      : undefined,
    amount_capturable: status === "requires_capture" ? amount : undefined,
    currency,
    payment_method: PAYMENT_METHODS[i % PAYMENT_METHODS.length],
    connector: CONNECTORS[i % CONNECTORS.length],
    error_code: status === "failed" ? "card_declined" : undefined,
    error_message:
      status === "failed" ? "The card was declined by the issuer." : undefined,
    confirmed_at: captured ? at(1 + (i % 28), 11) : undefined,
    captured_at: captured ? at(1 + (i % 28), 12) : undefined,
    cancelled_at: status === "cancelled" ? at(1 + (i % 28), 13) : undefined,
    created_at: at(1 + (i % 28), 10),
    updated_at: at(1 + (i % 28), 13),
  };
});

export const refunds: Refund[] = Array.from({ length: 14 }, (_, i) => {
  const payment = payments[(i * 2) % payments.length];
  return {
    uid: `99999999-9999-4999-8999-${pad(i)}`,
    refund_id: `ref_${String(i).padStart(20, "0")}`,
    payment_id: payment.payment_id,
    status: REFUND_STATUSES[i % REFUND_STATUSES.length],
    refund_type: i % 3 === 0 ? "scheduled" : "instant",
    amount: Math.floor(payment.amount / (1 + (i % 2))),
    currency: payment.currency,
    reason: i % 2 === 0 ? "Customer request" : undefined,
    connector: payment.connector,
    error_code: undefined,
    error_message: undefined,
    created_at: at(2 + (i % 26), 15),
    updated_at: undefined,
  };
});

const SETTLEMENT_BUCKETS = [
  { currency: "USD", method: "card" },
  { currency: "USD", method: "wallet" },
  { currency: "KRW", method: "card" },
];

export const settlements: Settlement[] = Array.from({ length: 24 }, (_, i) => {
  const day = 1 + Math.floor(i / SETTLEMENT_BUCKETS.length);
  const bucket = SETTLEMENT_BUCKETS[i % SETTLEMENT_BUCKETS.length];
  const date = `2026-09-${String(day).padStart(2, "0")}`;
  const unit = bucket.currency === "KRW" ? 1_000 : 100;
  const details: SettlementDetailLine[] = Array.from(
    { length: 3 + (i % 3) },
    (_, j) => {
      const isRefund = j === 2;
      return {
        uid: `aaaaaaaa-aaaa-4aaa-8aaa-${pad(i * 10 + j)}`,
        source_type: isRefund ? "REFUND" : "PAYMENT",
        source_id: isRefund
          ? `ref_${String(i * 10 + j).padStart(20, "0")}`
          : `pay_${String(i * 10 + j).padStart(20, "0")}`,
        payment_id: `pay_${String(i * 10 + (isRefund ? 0 : j)).padStart(20, "0")}`,
        amount: (isRefund ? -1 : 1) * unit * (12 + j * 3 + i),
        currency: bucket.currency,
        payment_method: bucket.method,
        occurred_at: at(day, 8 + j),
      };
    },
  );
  const gross = details
    .filter((d) => d.amount > 0)
    .reduce((sum, d) => sum + d.amount, 0);
  const refunded = details
    .filter((d) => d.amount < 0)
    .reduce((sum, d) => sum - d.amount, 0);
  const fee = Math.round(gross * 0.029);
  return {
    uid: `bbbbbbbb-bbbb-4bbb-8bbb-${pad(i)}`,
    settlement_date: date,
    currency: bucket.currency,
    payment_method: bucket.method,
    gross_amount: gross,
    refund_amount: refunded,
    fee_amount: fee,
    net_amount: gross - refunded - fee,
    payment_count: details.filter((d) => d.source_type === "PAYMENT").length,
    refund_count: details.filter((d) => d.source_type === "REFUND").length,
    status: day === 8 ? "OPEN" : i % 7 === 4 ? "DISCREPANCY" : "CONFIRMED",
    details,
    created_at: at(day + 1, 2),
    updated_at: day === 8 ? undefined : at(day + 1, 3),
  };
});

export const warehouses: Warehouse[] = [
  {
    uid: "cccccccc-cccc-4ccc-8ccc-000000000001",
    name: "Incheon Central",
    address: "12 Gonghang-ro, Jung-gu, Incheon",
    capacity: 120_000,
    created_at: "2026-01-05T00:00:00Z",
    updated_at: "2026-08-01T00:00:00Z",
  },
  {
    uid: "cccccccc-cccc-4ccc-8ccc-000000000002",
    name: "Busan Port",
    address: "45 Chungjang-daero, Busan",
    capacity: 80_000,
    created_at: "2026-02-10T00:00:00Z",
    updated_at: undefined,
  },
  {
    uid: "cccccccc-cccc-4ccc-8ccc-000000000003",
    name: "Oakland Overflow",
    address: undefined,
    capacity: 25_000,
    created_at: "2026-04-20T00:00:00Z",
    updated_at: undefined,
  },
];

export const stocks: Stock[] = Array.from({ length: 28 }, (_, i) => {
  const held = warehouses.filter((_, w) => (i + w) % 3 !== 2 || w === i % 3);
  const items = held.map((w, w2) => ({
    warehouse_uid: w.uid,
    warehouse_name: w.name,
    quantity: ((i + 1) * 17 + w2 * 11) % 240,
  }));
  return {
    uid: `dddddddd-dddd-4ddd-8ddd-${pad(i)}`,
    edition_uid: `eeeeeeee-eeee-4eee-8eee-${pad(i)}`,
    total_quantity: items.reduce((sum, item) => sum + item.quantity, 0),
    warehouses: items,
  };
});
