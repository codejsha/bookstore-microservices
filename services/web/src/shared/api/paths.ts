const V1 = "/api/v1";

export const paths = {
  works: {
    list: `${V1}/works`,
    search: `${V1}/works/search`,
    detail: (uid: string) => `${V1}/works/${uid}`,
    editions: (uid: string) => `${V1}/works/${uid}/editions`,
  },
  editions: {
    list: `${V1}/editions`,
    detail: (uid: string) => `${V1}/editions/${uid}`,
  },
  authors: {
    list: `${V1}/authors`,
    detail: (uid: string) => `${V1}/authors/${uid}`,
  },
  publishers: {
    list: `${V1}/publishers`,
    detail: (uid: string) => `${V1}/publishers/${uid}`,
  },
  subjects: {
    list: `${V1}/subjects`,
  },
  reviews: {
    byBook: (bookUid: string) => `${V1}/books/${bookUid}/reviews`,
  },
  cart: {
    root: `${V1}/cart`,
    items: `${V1}/cart/items`,
    item: (uid: string) => `${V1}/cart/items/${uid}`,
    checkout: `${V1}/cart/checkout`,
  },
  orders: {
    list: `${V1}/orders`,
    detail: (uid: string) => `${V1}/orders/${uid}`,
    cancel: (uid: string) => `${V1}/orders/${uid}/cancel`,
  },
  customers: {
    detail: (uid: string) => `${V1}/customers/${uid}`,
    reviews: (uid: string) => `${V1}/customers/${uid}/reviews`,
    wishlist: (uid: string) => `${V1}/customers/${uid}/wishlist`,
    wishlistAdd: (uid: string) => `${V1}/customers/${uid}/wishlist/add`,
    wishlistRemove: (uid: string) => `${V1}/customers/${uid}/wishlist/remove`,
    points: (uid: string) => `${V1}/customers/${uid}/points`,
    pointsHistory: (uid: string) => `${V1}/customers/${uid}/points/history`,
    pointsSpend: (uid: string) => `${V1}/customers/${uid}/points/spend`,
    pointsEarn: (uid: string) => `${V1}/customers/${uid}/points/earn`,
  },
} as const;
