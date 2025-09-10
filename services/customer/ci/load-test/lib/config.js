export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN || '';
const CUSTOMER_UID = __ENV.CUSTOMER_UID || '';
const BOOK_UID = __ENV.BOOK_UID || '';

export function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (TOKEN) h.Authorization = `Bearer ${TOKEN}`;
  return h;
}

export const READ_ENDPOINTS = [
  '/api/v1/customers?size=20',
  ...(CUSTOMER_UID
    ? [
        `/api/v1/customers/${CUSTOMER_UID}`,
        `/api/v1/customers/${CUSTOMER_UID}/wishlist`,
        `/api/v1/customers/${CUSTOMER_UID}/points`,
        `/api/v1/customers/${CUSTOMER_UID}/points/history`,
        `/api/v1/customers/${CUSTOMER_UID}/orders?size=20`,
        `/api/v1/customers/${CUSTOMER_UID}/reviews?size=20`,
      ]
    : []),
  ...(BOOK_UID ? [`/api/v1/books/${BOOK_UID}/reviews?size=20`] : []),
];
