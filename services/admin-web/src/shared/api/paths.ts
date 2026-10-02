const V1 = "/api/v1";

export const paths = {
  admin: {
    me: `${V1}/admin/me`,
    dashboard: `${V1}/admin/dashboard`,
  },
  risk: {
    list: `${V1}/admin/risk`,
    entry: (uid: string) => `${V1}/admin/risk/${uid}`,
  },
  catalog: {
    works: `${V1}/admin/catalog/works`,
    work: (uid: string) => `${V1}/admin/catalog/works/${uid}`,
    authors: `${V1}/admin/catalog/authors`,
    subjects: `${V1}/admin/catalog/subjects`,
  },
  users: {
    list: `${V1}/admin/users`,
    user: (uid: string) => `${V1}/admin/users/${uid}`,
  },
  orders: {
    list: `${V1}/admin/orders`,
    order: (uid: string) => `${V1}/admin/orders/${uid}`,
  },
  payments: {
    list: `${V1}/admin/payments`,
    payment: (uid: string) => `${V1}/admin/payments/${uid}`,
  },
  refunds: {
    list: `${V1}/admin/refunds`,
    refund: (uid: string) => `${V1}/admin/refunds/${uid}`,
  },
  settlements: {
    list: `${V1}/admin/settlements`,
    settlement: (uid: string) => `${V1}/admin/settlements/${uid}`,
  },
  inventory: {
    warehouses: `${V1}/admin/inventory/warehouses`,
    warehouse: (uid: string) => `${V1}/admin/inventory/warehouses/${uid}`,
    stocks: `${V1}/admin/inventory/stocks`,
  },
} as const;
