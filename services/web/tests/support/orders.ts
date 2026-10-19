import type { MockData } from "../../src/mocks/data";

type MockOrder = MockData["orders"][number];

export function buildOrders(count: number, template: MockOrder): MockOrder[] {
  return Array.from({ length: count }, (_, index) => {
    const n = index + 1;
    return {
      ...template,
      uid: `ord-gen-${n}`,
      order_number: `ORD-GEN-${String(n).padStart(4, "0")}`,
      idempotency_key: `idem-gen-${n}`,
    };
  });
}
