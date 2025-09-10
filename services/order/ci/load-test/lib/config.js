export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN || '';
const ORDER_UID = __ENV.ORDER_UID || '';

export function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (TOKEN) h.Authorization = `Bearer ${TOKEN}`;
  return h;
}

export const READ_ENDPOINTS = [
  '/api/v1/cart',
  '/api/v1/orders?size=20',
  ...(ORDER_UID ? [`/api/v1/orders/${ORDER_UID}`] : []),
];
