export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN || '';
const SHIPMENT_UID = __ENV.SHIPMENT_UID || '';

export function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (TOKEN) h.Authorization = `Bearer ${TOKEN}`;
  return h;
}

export const READ_ENDPOINTS = [
  '/api/v1/shipments?size=20',
  '/api/v1/carriers?size=20',
  '/api/v1/freights?size=20',
  '/api/v1/stats/dashboard',
  ...(SHIPMENT_UID
    ? [
        `/api/v1/shipments/${SHIPMENT_UID}`,
        `/api/v1/shipments/${SHIPMENT_UID}/tracking`,
        `/api/v1/freights/shipment/${SHIPMENT_UID}`,
      ]
    : []),
];
