import http from 'k6/http';
import { check, sleep } from 'k6';
import { BASE_URL, READ_ENDPOINTS, headers } from './lib/config.js';

export const options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '2m', target: 50 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(95)<800'],
  },
};

export default function () {
  const path = READ_ENDPOINTS[Math.floor(Math.random() * READ_ENDPOINTS.length)];
  const res = http.get(`${BASE_URL}${path}`, { headers: headers() });
  check(res, { 'status is 2xx': (r) => r.status >= 200 && r.status < 300 });
  sleep(0.5);
}
