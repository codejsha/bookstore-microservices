import http from 'k6/http';
import { check, sleep } from 'k6';
import { BASE_URL, READ_ENDPOINTS, headers } from './lib/config.js';

export const options = {
  vus: 1,
  duration: '30s',
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<500'],
  },
};

export default function () {
  for (const path of READ_ENDPOINTS) {
    const res = http.get(`${BASE_URL}${path}`, { headers: headers() });
    check(res, { 'status is 2xx': (r) => r.status >= 200 && r.status < 300 });
  }
  sleep(1);
}
