export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN || '';
const SETTLEMENT_UID = __ENV.SETTLEMENT_UID || '';

export function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (TOKEN) h.Authorization = `Bearer ${TOKEN}`;
  return h;
}

export const READ_ENDPOINTS = [
  '/api/v1/settlements?size=20',
  ...(SETTLEMENT_UID ? [`/api/v1/settlements/${SETTLEMENT_UID}`] : []),
];
