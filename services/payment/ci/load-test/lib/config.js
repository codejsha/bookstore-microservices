export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN || '';
const CUSTOMER_UID = __ENV.CUSTOMER_UID || '';
const PAYMENT_UID = __ENV.PAYMENT_UID || '';

export function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (TOKEN) h.Authorization = `Bearer ${TOKEN}`;
  return h;
}

export const READ_ENDPOINTS = [
  '/api/v1/customers?size=20',
  '/api/v1/payments?size=20',
  '/api/v1/refunds?size=20',
  '/api/v1/mandates?size=20',
  ...(CUSTOMER_UID
    ? [
        `/api/v1/customers/${CUSTOMER_UID}`,
        `/api/v1/customers/${CUSTOMER_UID}/payment-methods`,
      ]
    : []),
  ...(PAYMENT_UID
    ? [
        `/api/v1/payments/${PAYMENT_UID}`,
        `/api/v1/payments/${PAYMENT_UID}/attempts`,
      ]
    : []),
];
